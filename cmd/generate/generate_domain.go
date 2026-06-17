package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GenerateDomain(domainName string) error {
	if strings.Contains(domainName, "/") || strings.Contains(domainName, "\\") {
		return fmt.Errorf("invalid domain name: cannot contain path separators")
	}

	domain := cleanDomainName(domainName)

	if err := createDomainDirectories(domain); err != nil {
		return err
	}

	module := ""
	if ReadModuleName() != "" {
		module = ReadModuleName()
	}

	data := GeneratorData{
		Package: domain,
		Domain:  toSingularPascalCase(domain),
		Module:  module,
	}

	generateDomainFiles(domain, &data)
	generateModelsFile(domain, &data)

	err := AppendRouteToContainer(domain, toSingularPascalCase(domain), module)
	if err != nil {
		return fmt.Errorf("could not append routes to container: %v", err)
	}

	return nil
}

func cleanDomainName(name string) string {
	domain := strings.ToLower(name)
	if strings.HasSuffix(domain, "_") {
		domain = strings.ReplaceAll(domain, "_", "")
	}
	return domain
}

func createDomainDirectories(domain string) error {
	dir := filepath.Join("internal", domain)
	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		return DomainErr
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("cannot proceed creating domain folder %v", err)
	}

	if err := os.MkdirAll("internal/shared/models", 0755); err != nil {
		return fmt.Errorf("cannot proceed creating models folder %v", err)
	}

	return nil
}

func toSingularPascalCase(domain string) string {
	words := strings.Split(domain, "_")
	for i := range words {
		if len(words[i]) > 0 {
			words[i] = strings.ToUpper(words[i][:1]) + words[i][1:]
		}
	}
	result := strings.Join(words, "")

	if strings.HasSuffix(result, "ies") {
		return result[:len(result)-3] + "y"
	}
	if before, ok := strings.CutSuffix(result, "s"); ok {
		return before
	}
	return result
}

func generateDomainFiles(domain string, data *GeneratorData) {
	templates := map[string]string{
		"templates/test_handler.tmpl": fmt.Sprintf("%s_handler_test.go", domain),
		"templates/test_service.tmpl": fmt.Sprintf("%s_service_test.go", domain),
		"templates/mock.tmpl":         "mock.go",
		"templates/handler.tmpl":      "handler.go",
		"templates/service.tmpl":      "service.go",
		"templates/repository.tmpl":   "repository.go",
		"templates/types.tmpl":        "types.go",
	}

	for tmpPath, outputName := range templates {
		GenerateAndParse(domain, "internal", outputName, tmpPath, data)
	}
}

func generateModelsFile(domain string, data *GeneratorData) {
	modelPath := fmt.Sprintf("internal/shared/models/%s.go", domain)
	if _, err := os.Stat(modelPath); err != nil {
		err = GenerateAndParse("", "internal/shared/models", fmt.Sprintf("%s.go", domain), "templates/models.tmpl", data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Models file not created: %v\n", err)
		}
	}
}
