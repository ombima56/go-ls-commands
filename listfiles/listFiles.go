package listfiles

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"go-ls-commands/sorting"
)

func ListFiles(path string, longFormat bool, allFiles bool, recursive bool, timeSort bool, reverseSort bool, isFirst bool) {
	// Only show the ".: " header if recursive flag is set
	if isFirst && recursive {
		fmt.Println(path + ":")
	}

	// List current directory contents
	serveDir(path, longFormat, allFiles, timeSort, reverseSort)

	if recursive && len(path) > 1 {
		if strings.HasSuffix(path, "/") {
			path = strings.Trim(path, "/")
		}
	}

	if recursive {
		// Read directory contents
		files, err := os.ReadDir(path)
		if err != nil {
			fmt.Printf("cannot read directory '%s': %v\n", path, err)
			return
		}
		fileinfos := make([]fs.FileInfo, 0)
		for _, v := range files {
			fileinfo, _ := v.Info()
			fileinfos = append(fileinfos, fileinfo)
		}
		sorting.BubbleSortLowercaseFirst(fileinfos)
		if timeSort {
			sorting.SortTime(fileinfos)
		}
		if reverseSort {
			sorting.SortReverse(fileinfos)
		}

		// Collect directories for further processing
		var dirs []string
		for _, file := range fileinfos {
			if file.IsDir() && (allFiles || !strings.HasPrefix(file.Name(), ".")) {
				dirs = append(dirs, file.Name())
			}
		}

		if allFiles {
			// add parent and current directory...
		}

		// Process each directory
		for _, dirName := range dirs {

			fullPath := path + "/" + dirName
			// Convert absolute path to relative path for display

			fmt.Printf("\n%s:\n", fullPath)

			// Recursively list files in the subdirectory
			ListFiles(fullPath, longFormat, allFiles, recursive, timeSort, reverseSort, false)
		}
	}
}
