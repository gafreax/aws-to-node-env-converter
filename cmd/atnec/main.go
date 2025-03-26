package main

import (
	"fmt"
	"os"

	"aws-to-node-env-converter/internal/transformer"
)

func main() {
	// Check if correct number of arguments is provided
	if len(os.Args) != 3 {
		fmt.Println("Usage: atnec <input-file> <output-file>")
		os.Exit(1)
	}

	// Get input and output file paths
	inputPath := os.Args[1]
	outputPath := os.Args[2]

	// Transform the env file
	if err := transformer.TransformEnvFile(inputPath, outputPath); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully transformed %s to %s\n", inputPath, outputPath)
}
