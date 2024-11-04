package listfiles

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ombima56/go-ls-commands/listfiles"
)

func TestValidateFlags(t *testing.T) {
	tests := []struct {
		args        []string
		expectedErr bool
		expected    [5]bool // Represents longFlag, allFlag, recursiveFlag, timeFlag, reverseFlag
	}{
		{[]string{"-l"}, false, [5]bool{true, false, false, false, false}},
		{[]string{"-a"}, false, [5]bool{false, true, false, false, false}},
		{[]string{"-R"}, false, [5]bool{false, false, true, false, false}},
		{[]string{"-t"}, false, [5]bool{false, false, false, true, false}},
		{[]string{"-r"}, false, [5]bool{false, false, false, false, true}},
		{[]string{"--long"}, false, [5]bool{true, false, false, false, false}},
		{[]string{"--invalid"}, true, [5]bool{false, false, false, false, false}},
	}

	for _, test := range tests {
		longFlag, allFlag, recursiveFlag, timeFlag, reverseFlag, err := listfiles.ValidateFlags(test.args)
		if (err != nil) != test.expectedErr {
			t.Errorf("ValidateFlags(%v) error = %v, expectedErr = %v", test.args, err, test.expectedErr)
		}
		actual := [5]bool{longFlag, allFlag, recursiveFlag, timeFlag, reverseFlag}
		if actual != test.expected {
			t.Errorf("ValidateFlags(%v) = %v, want %v", test.args, actual, test.expected)
		}
	}
}
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

func TestListFiles(t *testing.T) {
	// Create temporary directory and files
	dir := t.TempDir()
	file1 := filepath.Join(dir, "file1.txt")
	os.WriteFile(file1, []byte("test file 1"), 0644)
	subDir := filepath.Join(dir, "subdir")
	os.Mkdir(subDir, 0755)
	file2 := filepath.Join(subDir, "file2.txt")
	os.WriteFile(file2, []byte("test file 2"), 0644)

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
