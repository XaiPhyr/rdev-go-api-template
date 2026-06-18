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

		for _, field := range structType.Fields.List {
			if field.Tag == nil || len(field.Names) == 0 {
				continue
			}

			tagVal := strings.Trim(field.Tag.Value, "`")
			structTag := reflect.StructTag(tagVal)
			bunTag := structTag.Get("bun")

			if bunTag == "" || bunTag == "-" {
				continue
			}

			var goType string
			switch t := field.Type.(type) {
			case *ast.Ident:
				goType = t.Name
			case *ast.SelectorExpr:
				if ident, ok := t.X.(*ast.Ident); ok {
					goType = ident.Name + "." + t.Sel.Name
				}
			case *ast.StarExpr:
				if selector, ok := t.X.(*ast.SelectorExpr); ok {
					if ident, ok := selector.X.(*ast.Ident); ok {
						goType = "*" + ident.Name + "." + selector.Sel.Name
					}
				} else if ident, ok := t.X.(*ast.Ident); ok {
					goType = "*" + ident.Name
				}
			}

			col := parseBunTagToSQL(field.Names[0].Name, goType, bunTag)
			columns = append(columns, col)
		}
		return false
	})

	return tableName, columns, nil
}

func parseBunTagToSQL(fieldName, goType, bunTag string) SQLColumn {
	parts := strings.Split(bunTag, ",")
	colName := strings.TrimSpace(parts[0])

	if colName == "" {
		for _, part := range parts {
			if strings.TrimSpace(part) == "soft_delete" {
				colName = "deleted_at"
			}
		}
	}

	if colName == "" {
		colName = strings.ToLower(fieldName)
	}

	var sqlType string
	switch goType {
	case "int", "int64":
		sqlType = "BIGINT"
	case "int32":
		sqlType = "INTEGER"
	case "float64", "float32":
		sqlType = "NUMERIC(10, 6)"
	case "time.Time", "*time.Time":
		sqlType = "TIMESTAMP WITH TIME ZONE"
	case "bool":
		sqlType = "BOOLEAN"
	default:
		sqlType = "VARCHAR(255)"
	}

	var constraints []string

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if after, ok := strings.CutPrefix(part, "type:"); ok {
			sqlType = strings.ToUpper(after)
			continue
		}

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
	}

	fullDefinition := sqlType
	if len(constraints) > 0 {
		fullDefinition = fmt.Sprintf("%s %s", sqlType, strings.Join(constraints, " "))
	}

	return SQLColumn{Name: colName, Type: fullDefinition}
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
