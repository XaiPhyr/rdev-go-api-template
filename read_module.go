package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ReadModuleName() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return ""
	}

	for {
		goModPath := filepath.Join(dir, "go.mod")

		if _, err := os.Stat(goModPath); err == nil {
			return parseGoMod(goModPath)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	fmt.Println("Error: go.mod file not found in current directory or any parent directories")
	return ""
}

func parseGoMod(path string) string {
	// #nosec G304 -- Path is safe; internally generated via walking up working directory for go.mod
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return ""
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") && !strings.HasPrefix(line, "//") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				return fields[1]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Scanner error:", err)
	}

	return ""
}
