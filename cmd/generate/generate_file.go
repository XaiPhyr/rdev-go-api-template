package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GenerateFile(domain, files string) error {
	if err := validateInput(domain, files); err != nil {
		return err
	}

	if err := ensureDomainDir(domain); err != nil {
		return err
	}

	for file := range strings.SplitSeq(files, ",") {
		if err := processFile(domain, file); err != nil {
			return err
		}
	}

	return nil
}

func validateInput(domain, files string) error {
	if strings.Contains(files, ".go") || strings.Contains(domain, "/") || strings.Contains(domain, "\\") {
		return fmt.Errorf("invalid file name: cannot contain file extension or path separators")
	}
	return nil
}

func ensureDomainDir(domain string) error {
	dir := filepath.Join("internal", strings.ToLower(domain))
	if _, err := os.Stat(dir); err != nil {
		if err := os.MkdirAll(dir, 0750); err != nil {
			return fmt.Errorf("cannot proceed creating domain folder %v", err)
		}
	}
	return nil
}

func processFile(domain, file string) error {
	targetPath := fmt.Sprintf("internal/%s/%s.go", domain, file)
	if _, err := os.Stat(targetPath); err == nil {
		return nil
	}

	cap := strings.ToUpper(domain[0:1])
	capitalizedDomain := fmt.Sprintf("%s%s", cap, domain[1:])

	module := ReadModuleName()

	data := GeneratorData{Package: domain, Domain: capitalizedDomain, Module: module}

	return generateTargetFile(domain, file, &data)
}

func generateTargetFile(domain, file string, data *GeneratorData) error {
	switch file {
	case "handler", "service", "repository", "types":
		return GenerateAndParse(domain, "internal", fmt.Sprintf("%s.go", file), fmt.Sprintf("templates/%s.tmpl", file), data)
	case "models":
		return generateModelsLayout(domain, data)
	case "tests":
		return generateTestLayouts(domain, data)
	default:
		return nil
	}
}

func generateModelsLayout(domain string, data *GeneratorData) error {
	if err := ensureModelsDir(); err != nil {
		return err
	}
	return GenerateAndParse("", "internal/shared/models", fmt.Sprintf("%s.go", domain), "templates/models.tmpl", data)
}

func generateTestLayouts(domain string, data *GeneratorData) error {
	if err := GenerateAndParse(domain, "internal", "mock.go", "templates/mock.tmpl", data); err != nil {
		return err
	}
	if err := GenerateAndParse(domain, "internal", fmt.Sprintf("%s_handler_test.go", domain), "templates/test_handler.tmpl", data); err != nil {
		return err
	}
	return GenerateAndParse(domain, "internal", fmt.Sprintf("%s_service_test.go", domain), "templates/test_service.tmpl", data)
}

func ensureModelsDir() error {
	modelsDir := "internal/shared/models"
	if _, err := os.Stat(modelsDir); err != nil {
		if err := os.MkdirAll(modelsDir, 0750); err != nil {
			return fmt.Errorf("cannot proceed creating models folder %v", err)
		}
	}
	return nil
}
