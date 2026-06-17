package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//go:embed cmd internal scripts compose.yaml config.sample.yaml go.mod .gitignore .github Dockerfile
var projectTemplates embed.FS

//go:embed templates/*
var tmpFiles embed.FS

// Regex matching alphanumeric characters, hyphens, underscores, and single periods
var validDirNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)

const oldModule = "github.com/XaiPhyr/rdev-go-api-template"

var ErrDomainExists = errors.New("domain Already Exists")
var ErrFileExists = errors.New("file Already Exists")

type GeneratorData struct {
	Package string
	Domain  string
	Module  string
}

func main() {
	if len(os.Args) < 2 {
		printGlobalUsage()
		os.Exit(1)
	}

	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
	domainFlag := generateCmd.String("d", "", "Domain name")
	fileFlag := generateCmd.String("f", "", "File to generate (handler|service|repository|types|models)")

	switch os.Args[1] {
	case "init":
		if len(os.Args) < 3 {
			fmt.Println("❌ Error: Missing module path argument.")
			fmt.Println("💡 Usage: rdev-go-api-template init <new-module-path>")
			os.Exit(1)
		}
		runProjectScaffolder(os.Args[2])

	case "generate":
		err := generateCmd.Parse(os.Args[2:])
		if err != nil {
			os.Exit(1)
		}

		if err := validateGeneratorFlags(*domainFlag, *fileFlag); err != nil {
			fmt.Fprintln(os.Stderr, err)
			generateCmd.Usage()
			os.Exit(1)
		}

		if err := runCodeGenerator(*domainFlag, *fileFlag); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Generation Error: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("❌ Error: Unknown sub-command '%s'\n", os.Args[1])
		printGlobalUsage()
		os.Exit(1)
	}
}

func validateGeneratorFlags(domain, file string) error {
	if file != "" && domain == "" {
		return fmt.Errorf("error: domain name flag (-f) must have flag (-d)")
	}
	if domain == "" {
		return fmt.Errorf("error: domain name flag (-d) is required")
	}
	return nil
}

func runCodeGenerator(domain, file string) error {
	if file != "" {
		return GenerateFile(domain, file)
	}
	return GenerateDomain(domain)
}

func runProjectScaffolder(newModule string) {
	targetDir, absTargetDir, err := validateScaffoldPaths(newModule)
	if err != nil {
		fmt.Printf("❌ Input Validation Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("🚀 Scaffolding new project into directory: %s...\n", absTargetDir)

	walker := createScaffoldWalker(targetDir, absTargetDir, newModule)

	if err := fs.WalkDir(projectTemplates, ".", walker); err != nil {
		fmt.Printf("❌ Failed to scaffold project: %v\n", err)
		os.Exit(1)
	}

	printScaffoldSuccess(targetDir, absTargetDir)
}

func validateScaffoldPaths(newModule string) (string, string, error) {
	parts := strings.Split(newModule, "/")
	unsafeTargetDir := parts[len(parts)-1]
	targetDir := filepath.Clean(unsafeTargetDir)

	if targetDir == "." || targetDir == ".." || strings.Contains(targetDir, "..") {
		return "", "", fmt.Errorf("invalid module name or path traversal attempt detected")
	}

	if !validDirNameRegex.MatchString(targetDir) {
		return "", "", fmt.Errorf("invalid project folder name '%s': special characters or spaces are not allowed (only alphanumeric, hyphens, and underscores are permitted)", targetDir)
	}

	absTargetDir, err := filepath.Abs(targetDir)
	if err != nil {
		return "", "", fmt.Errorf("error resolving absolute directory path: %w", err)
	}

	if _, err := os.Stat(absTargetDir); !os.IsNotExist(err) {
		return "", "", fmt.Errorf("folder './%s' already exists", targetDir)
	}

	return targetDir, absTargetDir, nil
}

func createScaffoldWalker(targetDir, absTargetDir, newModule string) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == "." || path == ".git" || strings.HasPrefix(path, ".git/") || path == "main.go" {
			return nil
		}

		relDestPath := filepath.Join(targetDir, path)
		absDestPath, err := filepath.Abs(relDestPath)
		if err != nil {
			return fmt.Errorf("error resolving absolute path for %s: %w", path, err)
		}

		if !strings.HasPrefix(absDestPath, absTargetDir) {
			return fmt.Errorf("security boundary breach: path traversal attempt blocked for %s", path)
		}

		if d.IsDir() {
			return os.MkdirAll(absDestPath, 0750)
		}

		return processAndWriteFile(absDestPath, path, newModule)
	}
}

func processAndWriteFile(absDestPath, embeddedPath, newModule string) error {
	parentDir := filepath.Dir(absDestPath)
	if err := os.MkdirAll(parentDir, 0750); err != nil {
		return fmt.Errorf("failed to create parent directory %s: %w", parentDir, err)
	}

	fileData, err := projectTemplates.ReadFile(embeddedPath)
	if err != nil {
		return err
	}

	if strings.Contains(embeddedPath, ".go") || embeddedPath == "go.mod" {
		fileData = []byte(strings.ReplaceAll(string(fileData), oldModule, newModule))
	}

	// #nosec G703 - Sanitized paths satisfy scanner taint flow engines safely
	return os.WriteFile(absDestPath, fileData, 0600)
}

func printScaffoldSuccess(targetDir, absTargetDir string) {
	fmt.Printf("\n✨ Success! Project initialized in %s\n", absTargetDir)
	fmt.Println("------------------------------------------------------------")
	fmt.Println("🎉 Your high-performance Go API architecture is ready!")
	fmt.Println("\n📌 Next Steps:")
	fmt.Printf("   📁  cd %s\n", targetDir)
	fmt.Println("   🛠️  go mod tidy")
	fmt.Println("------------------------------------------------------------")
}

func printGlobalUsage() {
	fmt.Println("Usage: rdev-go-api-template <command> [arguments]")
	fmt.Println("\nAvailable Commands:")
	fmt.Println("  init <module_name>              Scaffold a brand-new API repository structure")
	fmt.Println("  generate -d <domain>            Generate core architectural domain components inside an existing project")
	fmt.Println("  generate -d <domain> -f <file>  Generate a specific architectural component (handler|service|repository|types|models|tests)")
}
