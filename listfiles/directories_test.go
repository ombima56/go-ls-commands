package listfiles_test

import (
	"go-ls-commands/listfiles"
	"os"
	"testing"
	"time"
)

// createTempDir sets up a temporary directory with files and subdirectories for testing
func createTempDir(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()

	// Create files inside the temp directory
	file1 := tmpDir + "/file1.txt"
	file2 := tmpDir + "/file2.txt"
	subDir := tmpDir + "/subdir"

	_ = os.WriteFile(file1, []byte("File 1 content"), 0644)
	_ = os.WriteFile(file2, []byte("File 2 content"), 0644)
	_ = os.Mkdir(subDir, 0755)
	_ = os.WriteFile(subDir+"/file3.txt", []byte("File 3 content"), 0644)

	// Add a hidden file
	_ = os.WriteFile(tmpDir+"/.hidden.txt", []byte("Hidden file"), 0644)

	// Modify timestamps for sorting tests
	time.Sleep(1 * time.Second)
	_ = os.WriteFile(tmpDir+"/newer.txt", []byte("Newer file"), 0644)

	return tmpDir
}

func TestListFiles(t *testing.T) {
	tmpDir := createTempDir(t)
	opts := listfiles.Options{
		Recursive:   false,
		AllFiles:    false,
		LongFormat:  false,
		SortByTime:  false,
		ReverseSort: false,
	}

	t.Run("Normal file listing", func(t *testing.T) {
		listfiles.ListFiles(tmpDir, opts, true)
	})

	t.Run("Recursive listing", func(t *testing.T) {
		opts.Recursive = true
		listfiles.ListFiles(tmpDir, opts, true)
	})

	t.Run("Listing with hidden files", func(t *testing.T) {
		opts.Recursive = false
		opts.AllFiles = true
		listfiles.ListFiles(tmpDir, opts, true)
	})

	t.Run("Sorting by time", func(t *testing.T) {
		opts.SortByTime = true
		listfiles.ListFiles(tmpDir, opts, true)
	})

	t.Run("Reverse sorting", func(t *testing.T) {
		opts.SortByTime = false
		opts.ReverseSort = true
		listfiles.ListFiles(tmpDir, opts, true)
	})
}

func TestCalculateFileMetadata(t *testing.T) {
	tmpDir := createTempDir(t)

	files, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	fileInfos := make([]os.FileInfo, 0)
	for _, file := range files {
		info, _ := file.Info()
		fileInfos = append(fileInfos, info)
	}

	metadata := listfiles.CalculateFileMetadata(fileInfos)

	if metadata.MaxSize == 0 {
		t.Errorf("Expected a non-zero max size, got %d", metadata.MaxSize)
	}
}
