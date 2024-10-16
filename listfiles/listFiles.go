package listfiles

import (
	"fmt"
	"os"
	"sort"
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
	
	// Sort files based on the specified criteria
	if sortByTime {
		sort.Slice(files, func(i, j int) bool {
			if reverseOrder {
				return files[i].ModTime.After(files[j].ModTime)
			}
			return files[i].ModTime.Before(files[j].ModTime)
		})
	} else {
		sort.Slice(files, func(i, j int) bool {
			if reverseOrder {
				return files[i].Name > files[j].Name
			}
			return files[i].Name < files[j].Name
		})
	}

	// Display the files
	for _, file := range files {
		if showDetails {
			fileType := "FILE"
			if file.IsDir {
				fileType = "DIR"
			}
			fmt.Printf("%s\t%s\n", fileType, file.Name)
		} else {
			fmt.Println(file.Name)
		}
	}
	return nil
}
