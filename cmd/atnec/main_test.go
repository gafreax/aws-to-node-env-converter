package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCLIMainFunction(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := ioutil.TempDir("", "atnec-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Prepare test input file
	inputPath := filepath.Join(tempDir, "input.env")
	inputContent := "KEY1=VALUE1\n# Comment\nKEY2=VALUE2"
	err = ioutil.WriteFile(inputPath, []byte(inputContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Prepare output file path
	outputPath := filepath.Join(tempDir, "output.env")

	// Build the CLI binary
	binaryPath := filepath.Join(tempDir, "atnec")
	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/atnec")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}

	// Run the CLI with input and output paths
	cliCmd := exec.Command(binaryPath, inputPath, outputPath)
	var outBuf, errBuf bytes.Buffer
	cliCmd.Stdout = &outBuf
	cliCmd.Stderr = &errBuf

	err = cliCmd.Run()
	if err != nil {
		t.Fatalf("CLI command failed: %v\nStdout: %s\nStderr: %s",
			err, outBuf.String(), errBuf.String())
	}

	// Check output file contents
	outputContent, err := ioutil.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expectedOutput := "KEY1='VALUE1'\nKEY2='VALUE2'\n"
	if string(outputContent) != expectedOutput {
		t.Errorf("Unexpected output.\nExpected:\n%s\nGot:\n%s",
			expectedOutput, string(outputContent))
	}

	// Check stdout
	expectedStdout := "Successfully transformed " + inputPath + " to " + outputPath + "\n"
	if outBuf.String() != expectedStdout {
		t.Errorf("Unexpected stdout.\nExpected: %s\nGot: %s",
			expectedStdout, outBuf.String())
	}
}

func TestCLIErrorHandling(t *testing.T) {
	testCases := []struct {
		name          string
		args          []string
		expectedError bool
	}{
		{
			name:          "No arguments",
			args:          []string{},
			expectedError: true,
		},
		{
			name:          "Single argument",
			args:          []string{"input.env"},
			expectedError: true,
		},
		{
			name:          "Too many arguments",
			args:          []string{"input.env", "output.env", "extra.env"},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Build the CLI binary
			tempDir, err := ioutil.TempDir("", "atnec-error-test-")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tempDir)

			binaryPath := filepath.Join(tempDir, "atnec")
			cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/atnec")
			if err := cmd.Run(); err != nil {
				t.Fatalf("Failed to build CLI: %v", err)
			}

			// Prepare CLI command
			cliCmd := exec.Command(binaryPath, tc.args...)
			var outBuf, errBuf bytes.Buffer
			cliCmd.Stdout = &outBuf
			cliCmd.Stderr = &errBuf

			err = cliCmd.Run()
			if tc.expectedError && err == nil {
				t.Errorf("Expected error, but command succeeded")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check error output for incorrect argument count
			if tc.expectedError {
				expectedErrorOutput := "Usage: atnec <input-file> <output-file>\n"
				if outBuf.String() != expectedErrorOutput {
					t.Errorf("Unexpected error output.\nExpected: %s\nGot: %s",
						expectedErrorOutput, outBuf.String())
				}
			}
		})
	}
}
