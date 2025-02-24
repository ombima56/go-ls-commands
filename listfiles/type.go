package listfiles

import (
	"os"
)

// CustomFileInfo wraps os.FileInfo to override the Name() method
type CustomFileInfo struct {
	os.FileInfo
	name string
}

// Name overrides the original FileInfo's Name method
func (f CustomFileInfo) Name() string {
	return f.name
}

// FileMetadata holds the maximum field lengths for formatting
type FileMetadata struct {
	MaxFieldLengths map[string]int
	MaxSize         int64
}

// NewFileMetadata creates a new metadata struct with initialized map
func NewFileMetadata() FileMetadata {
	return FileMetadata{
		MaxFieldLengths: map[string]int{
			"permissions": 0,
			"links":       0,
			"owner":       0,
			"group":       0,
			"size":        0,
			"modTime":     0,
			"fileName":    0,
		},
		MaxSize: 0,
	}
}