package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

type SQLColumn struct {
	Name string
	Type string
}

func ParseBunModels(filePath string) (string, []SQLColumn, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse Go file: %w", err)
	}

	var tableName string
	var columns []SQLColumn

	ast.Inspect(node, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		tableName = strings.ToLower(typeSpec.Name.Name) + "s"
		columns = extractColumnsFromStruct(structType)
		return false
	})

	return tableName, columns, nil
}

func extractColumnsFromStruct(structType *ast.StructType) []SQLColumn {
	var columns []SQLColumn
	for _, field := range structType.Fields.List {
		if field.Tag == nil || len(field.Names) == 0 {
			continue
		}

		structTag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
		bunTag := structTag.Get("bun")
		if bunTag == "" || bunTag == "-" {
			continue
		}

		goType := determineGoTypeString(field.Type)
		col := parseBunTagToSQL(field.Names[0].Name, goType, bunTag)
		columns = append(columns, col)
	}
	return columns
}

func determineGoTypeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name + "." + t.Sel.Name
		}
	case *ast.StarExpr:
		if selector, ok := t.X.(*ast.SelectorExpr); ok {
			if ident, ok := selector.X.(*ast.Ident); ok {
				return "*" + ident.Name + "." + selector.Sel.Name
			}
		}
		if ident, ok := t.X.(*ast.Ident); ok {
			return "*" + ident.Name
		}
	}
	return ""
}

func parseBunTagToSQL(fieldName, goType, bunTag string) SQLColumn {
	parts := strings.Split(bunTag, ",")
	colName := resolveColumnName(fieldName, parts)
	sqlType := mapGoTypeToSQL(goType)

	var constraints []string
	sqlType, constraints = processTagConstraints(parts, sqlType, constraints)

	fullDefinition := sqlType
	if len(constraints) > 0 {
		fullDefinition = fmt.Sprintf("%s %s", sqlType, strings.Join(constraints, " "))
	}

	return SQLColumn{Name: colName, Type: fullDefinition}
}

func resolveColumnName(fieldName string, parts []string) string {
	colName := strings.TrimSpace(parts[0])
	if colName == "" {
		for _, part := range parts {
			if strings.TrimSpace(part) == "soft_delete" {
				return "deleted_at"
			}
		}
		return strings.ToLower(fieldName)
	}
	return colName
}

func mapGoTypeToSQL(goType string) string {
	switch goType {
	case "int", "int64":
		return "BIGINT"
	case "int32":
		return "INTEGER"
	case "float64", "float32":
		return "NUMERIC(10, 6)"
	case "time.Time", "*time.Time":
		return "TIMESTAMP WITH TIME ZONE"
	case "bool":
		return "BOOLEAN"
	default:
		return "VARCHAR(255)"
	}
}

func processTagConstraints(parts []string, sqlType string, constraints []string) (string, []string) {
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if after, ok := strings.CutPrefix(part, "type:"); ok {
			sqlType = strings.ToUpper(after)
			continue
		}
		sqlType, constraints = appendConstraint(part, sqlType, constraints)
	}
	return sqlType, constraints
}

func appendConstraint(part, sqlType string, constraints []string) (string, []string) {
	switch {
	case part == "pk":
		constraints = append(constraints, "PRIMARY KEY")
	case part == "autoincrement":
		switch sqlType {
		case "BIGINT":
			sqlType = "BIGSERIAL"
		case "INTEGER":
			sqlType = "SERIAL"
		}
	case part == "notnull":
		constraints = append(constraints, "NOT NULL")
	case part == "unique":
		constraints = append(constraints, "UNIQUE")
	case strings.HasPrefix(part, "default:"):
		defVal := strings.TrimPrefix(part, "default:")
		constraints = append(constraints, fmt.Sprintf("DEFAULT %s", defVal))
	}
	return sqlType, constraints
}

func BuildSQLMigration(tableName string, columns []SQLColumn) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "-- Migration Up: Create %s table\n", tableName)
	fmt.Fprintf(&sb, "CREATE TABLE IF NOT EXISTS %s (\n", tableName)

	var colDefs []string
	for _, col := range columns {
		colDefs = append(colDefs, fmt.Sprintf("    %s %s", col.Name, col.Type))
	}

	sb.WriteString(strings.Join(colDefs, ",\n"))
	sb.WriteString("\n);\n")
	return sb.String()
}
