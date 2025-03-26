package transformer

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// Helper function to create a temp file with given content
func createTempFile(t *testing.T, content string) string {
	t.Helper()
	tmpfile, err := ioutil.TempFile("", "test-env-*.env")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	return tmpfile.Name()
}

// Helper function to read file contents
func readFileContents(t *testing.T, path string) string {
	t.Helper()
	content, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

func TestTransformEnvFile(t *testing.T) {
	testCases := []struct {
		name           string
		inputContent   string
		expectedOutput string
	}{
		{
			name: "Basic transformation",
			inputContent: `KEY1=VALUE1
KEY2=VALUE2`,
			expectedOutput: "KEY1='VALUE1'\nKEY2='VALUE2'\n",
		},
		{
			name: "With comments and empty lines",
			inputContent: `# This is a comment
KEY1=VALUE1

# Another comment
KEY2=VALUE2`,
			expectedOutput: "KEY1='VALUE1'\nKEY2='VALUE2'\n",
		},
		{
			name: "With whitespace",
			inputContent: `  KEY1  =  VALUE1  
  KEY2 = VALUE2  `,
			expectedOutput: "KEY1='VALUE1'\nKEY2='VALUE2'\n",
		},
		{
			name: "Complex values",
			inputContent: `DB_CONNECTION=mysql://user:pass@host:port/database
API_KEY=abc123!@#`,
			expectedOutput: "DB_CONNECTION='mysql://user:pass@host:port/database'\nAPI_KEY='abc123!@#'\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temporary input file
			inputPath := createTempFile(t, tc.inputContent)
			defer os.Remove(inputPath)

			// Create temporary output file
			outputPath := filepath.Join(os.TempDir(), "output-test.env")
			defer os.Remove(outputPath)

			// Transform the file
			err := TransformEnvFile(inputPath, outputPath)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Read and check output file contents
			outputContent := readFileContents(t, outputPath)
			if outputContent != tc.expectedOutput {
				t.Errorf("Unexpected output.\nExpected:\n%s\nGot:\n%s",
					tc.expectedOutput, outputContent)
			}
		})
	}
}

func TestTransformEnvFileErrorHandling(t *testing.T) {
	t.Run("Non-existent input file", func(t *testing.T) {
		outputPath := filepath.Join(os.TempDir(), "output-test.env")
		err := TransformEnvFile("/path/to/non/existent/file", outputPath)
		if err == nil {
			t.Fatal("Expected an error for non-existent input file")
		}
	})

	t.Run("Invalid output path", func(t *testing.T) {
		inputPath := createTempFile(t, "KEY=VALUE")
		defer os.Remove(inputPath)

		err := TransformEnvFile(inputPath, "/invalid/path/output.env")
		if err == nil {
			t.Fatal("Expected an error for invalid output path")
		}
	})
}
