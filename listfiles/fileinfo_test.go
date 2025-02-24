package listfiles

import (
	"os"
	"strconv"
	"testing"
)

// TestUpdateFieldLengths ensures correct max field length calculations
func TestUpdateFieldLengths(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Get file info
	fileInfo, err := tmpFile.Stat()
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	// Initialize maxLengths map
	maxLengths := map[string]int{
		"permissions": 0, "links": 0, "owner": 0, "group": 0, "size": 0, "modTime": 0, "fileName": 0,
	}

	updateFieldLengths(fileInfo, maxLengths)

	// Validate updates
	if maxLengths["permissions"] < len(FileModeToString(fileInfo.Mode())) {
		t.Errorf("Permissions length not updated correctly")
	}
	if maxLengths["size"] < len(strconv.Itoa(int(fileInfo.Size()))) {
		t.Errorf("Size length not updated correctly")
	}
}
