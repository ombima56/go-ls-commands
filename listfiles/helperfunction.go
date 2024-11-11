package listfiles

import (
	"os"
	"time"
)

type FileInfo struct {
	Name    string
	IsDir   bool
	ModTime time.Time
}

type customFileInfo struct {
	os.FileInfo
	name string
}

// Helper function to create a custom FileInfo for . and ..
func (f customFileInfo) Name() string {
	return f.name
}
