package listfiles

import (
	"os"
	"strings"
	"time"
)

type FileInfo struct {
	Name    string
	IsDir   bool
	ModTime time.Time
}

func ListFiles(dir string, showDetails, recursive, includeHidden, reverseOrder, sortByTime bool) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var files []FileInfo

	// Collecting all file information
	for _, entry := range entries {
		if !includeHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}

		files = append(files, FileInfo{
			Name:    info.Name(),
			IsDir:   info.IsDir(),
			ModTime: info.ModTime(),
		})
	}
	return nil
}
