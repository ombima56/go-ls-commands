package sorting

import (
	"io/fs"
	"testing"
	"time"
)

// mockFileInfo implements fs.FileInfo interface for testing
type mockFileInfo struct {
	name    string
	modTime time.Time
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return 0 }
func (m mockFileInfo) Mode() fs.FileMode  { return 0 }
func (m mockFileInfo) ModTime() time.Time { return m.modTime }
func (m mockFileInfo) IsDir() bool        { return false }
func (m mockFileInfo) Sys() any          { return nil }

func TestSortTime(t *testing.T) {
	tests := []struct {
		name     string
		files    []fs.FileInfo
		expected []string // expected order of file names after sorting
	}{
		{
			name: "files in descending order",
			files: []fs.FileInfo{
				mockFileInfo{"file1.txt", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
				mockFileInfo{"file2.txt", time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)},
				mockFileInfo{"file3.txt", time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)},
			},
			expected: []string{"file3.txt", "file2.txt", "file1.txt"},
		},
		{
			name: "files in random order",
			files: []fs.FileInfo{
				mockFileInfo{"file2.txt", time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)},
				mockFileInfo{"file1.txt", time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)},
				mockFileInfo{"file3.txt", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			expected: []string{"file1.txt", "file2.txt", "file3.txt"},
		},
		{
			name: "files with same modification time",
			files: []fs.FileInfo{
				mockFileInfo{"file1.txt", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
				mockFileInfo{"file2.txt", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			expected: []string{"file1.txt", "file2.txt"},
		},
		{
			name:     "empty slice",
			files:    []fs.FileInfo{},
			expected: []string{},
		},
		{
			name: "single file",
			files: []fs.FileInfo{
				mockFileInfo{"file1.txt", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			expected: []string{"file1.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SortTime(tt.files)

			// Check if files are sorted correctly
			for i, expected := range tt.expected {
				if tt.files[i].Name() != expected {
					t.Errorf("SortTime() got file order %v at position %d, want %v",
						tt.files[i].Name(), i, expected)
				}
			}

			// Verify descending order of modification times
			for i := 1; i < len(tt.files); i++ {
				if tt.files[i-1].ModTime().Before(tt.files[i].ModTime()) {
					t.Errorf("SortTime() files not in descending order: %v before %v",
						tt.files[i-1].ModTime(), tt.files[i].ModTime())
				}
			}
		})
	}
}