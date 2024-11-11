package listfiles

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	filepaths "go-ls-commands/filepath"
	"go-ls-commands/sorting"
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

		// parentDir := filepath.Dir(dir)
		parentDir := filepaths.GetParentDir(dir)
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
	if timeSort {
		sorting.SortTime(fileInfos)
	} else {
		sorting.BubbleSortLowercaseFirst(fileInfos)
	}

	if reverseSort {
		sorting.SortReverse(fileInfos)
	}

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

	maxSize := GetMaxFileSize(fileInfos)

	for _, file := range fileInfos {
		if longFormat {
			PrintFileInfo(dir, file, maxSize)
		} else {
			PrintFileName(file)
		}
	}

	// Add newline if not in long format
	if !longFormat {
		fmt.Println()
	}
}
