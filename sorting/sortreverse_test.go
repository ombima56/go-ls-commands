package sorting

import (
	"io/fs"
	"testing"
	"time"
)

// MockFileInfo is a mock implementation of the fs.FileInfo interface for testing purposes.
type MockFileInfo struct {
	name    string
	size    int64
	modTime time.Time
	isDir   bool
}

func (m MockFileInfo) Name() string       { return m.name }
func (m MockFileInfo) Size() int64        { return m.size }
func (m MockFileInfo) Mode() fs.FileMode  { return 0 }
func (m MockFileInfo) ModTime() time.Time { return m.modTime }
func (m MockFileInfo) IsDir() bool        { return m.isDir }
func (m MockFileInfo) Sys() interface{}   { return nil }

func TestSortReverse(t *testing.T) {
	type args struct {
		files []fs.FileInfo
	}
	tests := []struct {
		name     string
		args     args
		expected []fs.FileInfo
	}{
		{
			name: "reverse sorted files",
			args: args{
				files: []fs.FileInfo{
					MockFileInfo{name: "file1"},
					MockFileInfo{name: "file2"},
					MockFileInfo{name: "file3"},
				},
			},
			expected: []fs.FileInfo{
				MockFileInfo{name: "file3"},
				MockFileInfo{name: "file2"},
				MockFileInfo{name: "file1"},
			},
		},
		{
			name: "empty file list",
			args: args{
				files: []fs.FileInfo{},
			},
			expected: []fs.FileInfo{},
		},
		{
			name: "single file",
			args: args{
				files: []fs.FileInfo{
					MockFileInfo{name: "file1"},
				},
			},
			expected: []fs.FileInfo{
				MockFileInfo{name: "file1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SortReverse(tt.args.files)
			// Check if the result matches the expected output
			for i, file := range tt.args.files {
				if file.Name() != tt.expected[i].Name() {
					t.Errorf("SortReverse() failed for test %s: got %v, want %v", tt.name, file.Name(), tt.expected[i].Name())
				}
			}
		})
	}
}
