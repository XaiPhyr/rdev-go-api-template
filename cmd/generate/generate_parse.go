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

	tmpl, err := template.New(outputName).Parse(string(tmplBytes))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", outputName, err)
	}

	// issue: os.Create(targetFilepath)
	// what?: not creating targetFilePath error failed to create specific path
	// how?: because targetFilePath only generates the domain and the file,
	// since we are generating it for internal, we need to make sure to join it using filepath.Join
	file, err := os.Create(filepath.Join(folder, targetFilePath))
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", targetFilePath, err)
	}

	err = tmpl.Execute(file, data)
	file.Close()

	if err != nil {
		return fmt.Errorf("failed to write template %s: %w", targetFilePath, err)
	}

	fmt.Printf(" -> Created: %s/%s\n", folder, outputName)

	return nil
}
