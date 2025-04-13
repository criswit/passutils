package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const defaultOutputDirectory = "password-export"
const outputFileName = "password-export.json"

func main() {
	// Define command-line flag for output directory
	outputDir := flag.String("outdir", defaultOutputDirectory, "output directory for exported passwords")
	flag.Parse()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	exporter := &PassExporter{baseDir: fmt.Sprintf("%s/.password-store", homeDir), data: make(map[string]interface{})}

	if err := exporter.Export(); err != nil {
		panic(exporter)
	}

	jsonData, err := json.MarshalIndent(exporter.data, "", "   ")
	if err != nil {
		panic(err)
	}

	// Determine output directory path
	var outputDirectoryAbsolutePath string
	if filepath.IsAbs(*outputDir) {
		// If user provided absolute path, use it directly
		outputDirectoryAbsolutePath = *outputDir
	} else {
		// If user provided relative path or default is used, append to home directory
		outputDirectoryAbsolutePath = filepath.Join(homeDir, *outputDir)
	}

	if err := createDirIfNotExists(outputDirectoryAbsolutePath); err != nil {
		panic(err)
	}

	outputFilePath := filepath.Join(outputDirectoryAbsolutePath, outputFileName)
	if err := os.WriteFile(outputFilePath, jsonData, 0600); err != nil {
		panic(err)
	}

	fmt.Printf("Passwords exported to: %s\n", outputFilePath)
}

func createDirIfNotExists(path string) error {
	// Check if directory already exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Directory does not exist, create it with permissions set to 0755
		// (owner can read/write/execute, others can read/execute)
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		fmt.Printf("Created directory: %s\n", path)
	} else if err != nil {
		// Some other error occurred while checking
		return fmt.Errorf("error checking directory: %w", err)
	} else {
		// Directory already exists
		fmt.Printf("Directory already exists: %s\n", path)
	}

	return nil
}
