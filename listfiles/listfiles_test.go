package listfiles_test

import (
	"go-ls-commands/listfiles"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListFiles(t *testing.T) {
	// Create temporary directory and files
	dir := t.TempDir()
	file1 := filepath.Join(dir, "file1.txt")
	os.WriteFile(file1, []byte("test file 1"), 0o644)
	subDir := filepath.Join(dir, "subdir")
	os.Mkdir(subDir, 0o755)
	file2 := filepath.Join(subDir, "file2.txt")
	os.WriteFile(file2, []byte("test file 2"), 0o644)

	// Capture output of ListFiles for testing
	output := getOutput(func() {
		listfiles.ListFiles(dir, true, true, true, false, false, true)
	})

	// Check if the output contains file1 and subDir as expected
	if !strings.Contains(output, "file1.txt") || !strings.Contains(output, "subdir") {
		t.Errorf("ListFiles output = %v; want to contain file1.txt and subdir", output)
	}

	// Test with different flag settings (e.g., reverseSort, timeSort)
	output = getOutput(func() {
		listfiles.ListFiles(dir, true, true, true, true, true, true)
	})

	// Check if it contains necessary information based on flags
	if !strings.Contains(output, "file1.txt") || !strings.Contains(output, "subdir") {
		t.Errorf("ListFiles output = %v; want to contain file1.txt and subdir", output)
	}
}
