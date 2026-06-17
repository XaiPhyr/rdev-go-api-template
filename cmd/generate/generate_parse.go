package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func GenerateAndParse(domain, folder, outputName, tmpPath string, data *GeneratorData) error {
	targetFilePath := filepath.Join(domain, outputName)

	tmplBytes, err := tmpFiles.ReadFile(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to read template source: %w", err)
	}

	tmpl, err := template.New(outputName).Parse(string(tmplBytes))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", outputName, err)
	}

	fullPath := filepath.Join(folder, targetFilePath)

	// #nosec G304 -- Safe path construction for local file generation scaffolding
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", targetFilePath, err)
	}

	defer func() { _ = file.Close() }()

	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("failed to write template %s: %w", targetFilePath, err)
	}

	fmt.Printf(" -> Created: %s/%s\n", folder, outputName)
	return nil
}
