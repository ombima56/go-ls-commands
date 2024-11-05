package listfiles

import (
	"os"
	"time"
)

const (
	Reset = "\033[0m"
	Blue  = "\033[34m" // Directory color
	Green = "\033[32m" // File color
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
