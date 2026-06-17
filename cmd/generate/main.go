package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"os"
)

//go:embed templates/*
var tmpFiles embed.FS

type GeneratorData struct {
	Package string
	Domain  string
	Module  string
}

var ErrDomainExists = errors.New("domain Already Exists")
var ErrFileExists = errors.New("file Already Exists")

func main() {
	domainFlag := flag.String("d", "", "Domain name")
	fileFlag := flag.String("f", "", "File to generate (handler|service|repository|types|models) (e.g. generate -d orders -f handler)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUsage of generate:\n")
		fmt.Fprintf(os.Stderr, "  generate -d <domain_name>\n")
		fmt.Fprintf(os.Stderr, "  generate -f <file>\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if err := validateFlags(*domainFlag, *fileFlag); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		os.Exit(1)
	}

	if err := runGenerator(*domainFlag, *fileFlag); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		if errors.Is(err, ErrFileExists) || errors.Is(err, ErrDomainExists) {
			os.Exit(1)
		}
		flag.Usage()
		os.Exit(1)
	}
}

func validateFlags(domain, file string) error {
	if len(os.Args) < 2 {
		return fmt.Errorf("error: missing arguments")
	}
	if file != "" && domain == "" {
		return fmt.Errorf("error: domain name flag (-f) must have flag (-d)")
	}
	if domain == "" {
		return fmt.Errorf("error: domain name flag (-d) is required")
	}
	return nil
}

func runGenerator(domain, file string) error {
	if file != "" {
		return GenerateFile(domain, file)
	}
	return GenerateDomain(domain)
}
