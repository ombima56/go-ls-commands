package listfiles_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	filepaths "go-ls-commands/filepath"
	"go-ls-commands/listfiles"
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
	file := filepaths.JoinPaths(t.TempDir(), "testfile.txt")
	os.WriteFile(file, []byte("test data"), 0o644)

	// Stat the file to get its FileInfo
	fileInfo, err := os.Stat(file)
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}

	
	maxSize := listfiles.GetMaxFileSize([]os.FileInfo{fileInfo})

	// Capture output of PrintFileInfo for testing
	output := getOutput(func() {
		listfiles.PrintFileInfo(file, fileInfo, maxSize)
	})

	// Check if the output contains file permissions, owner, size, etc.
	expectedName := "testfile.txt"
	if !strings.Contains(output, expectedName) {
		t.Errorf("PrintFileInfo output = %v; want to contain %v", output, expectedName)
	}
}
