package listfiles_test

import (
	"bytes"
	"my-ls/listfiles"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func getOutput(f func()) string {
	var outputedRes bytes.Buffer
	stdout := os.Stdout

	defer func() {
		os.Stdout = stdout
	}()

	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	outputedRes.ReadFrom(r)

	return outputedRes.String()
}

func TestPrintFileInfo(t *testing.T) {
	// Create a temporary file
	file := filepath.Join(t.TempDir(), "testfile.txt")
	os.WriteFile(file, []byte("test data"), 0644)

	// Stat the file to get its FileInfo
	fileInfo, err := os.Stat(file)
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}

	// Capture output of PrintFileInfo for testing
	output := getOutput(func() {
		listfiles.PrintFileInfo(fileInfo)
	})

	// Check if the output contains file permissions, owner, size, etc.
	expectedName := "testfile.txt"
	if !strings.Contains(output, expectedName) {
		t.Errorf("PrintFileInfo output = %v; want to contain %v", output, expectedName)
	}
}
