package listfiles

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
)

func serveDir(dir string, longFormat bool, allFiles bool, timeSort bool, reverseSort bool) {
	f, err := os.Open(dir)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer f.Close()

	// Read all files in the directory
	files, err := f.Readdir(0)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	var fileInfos []os.FileInfo

	// Add . and .. if allFiles is set
	if allFiles {
		curDirInfo, _ := os.Stat(dir)
		fileInfos = append(fileInfos, customFileInfo{curDirInfo, "."})

		parentDir := filepath.Dir(dir)
		parentDirInfo, _ := os.Stat(parentDir)
		fileInfos = append(fileInfos, customFileInfo{parentDirInfo, ".."})
	}

	// Add other files
	for _, file := range files {
		if !allFiles && strings.HasPrefix(file.Name(), ".") {
			continue // skip hidden files if -a is not set
		}
		fileInfos = append(fileInfos, file)
	}

	// Sort files based on the specified criteria
	sort.Slice(fileInfos, func(i, j int) bool {
		nameI := fileInfos[i].Name()
		nameJ := fileInfos[j].Name()

		// Always put . and .. first
		if nameI == "." || nameI == ".." {
			return nameJ != "." && nameJ != ".."
		}
		if nameJ == "." || nameJ == ".." {
			return false
		}

		// Determine sorting order based on timeSort flag
		if timeSort {
			if fileInfos[i].ModTime().After(fileInfos[j].ModTime()) {
				return !reverseSort // Return true if not reversed
			}
			return reverseSort // Return true if reversed
		}

		// Default name sorting
		if nameI < nameJ {
			return !reverseSort // Return true if not reversed
		}
		return reverseSort // Return true if reversed
	})

	// Print files
	if longFormat {
		var totalBlocks int64
		for _, file := range fileInfos {
			if stat, ok := file.Sys().(*syscall.Stat_t); ok {
				totalBlocks += int64(stat.Blocks)
			}
		}
		fmt.Printf("total %d\n", totalBlocks/2)
	}

	for _, file := range fileInfos {
		if longFormat {
			PrintFileInfo(file)
		} else {
			PrintFileName(file)
		}
	}

	// Add newline if not in long format
	if !longFormat {
		fmt.Println()
	}
}
