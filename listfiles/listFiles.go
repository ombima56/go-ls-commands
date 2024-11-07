package listfiles

import (
	"fmt"
	"os"
	"strings"
)

func ListFiles(path string, longFormat bool, allFiles bool, recursive bool, timeSort bool, reverseSort bool, isFirst bool) {
	// Only show the ".: " header if recursive flag is set
	if isFirst && recursive {
		fmt.Println(".:")
	}

	// List current directory contents
	serveDir(path, longFormat, allFiles, timeSort, reverseSort)

	if recursive {
		// Read directory contents
		files, err := os.ReadDir(path)
		if err != nil {
			fmt.Printf("cannot read directory '%s': %v\n", path, err)
			return
		}

		// Collect directories for further processing
		var dirs []string
		for _, file := range files {
			if file.IsDir() && (allFiles || !strings.HasPrefix(file.Name(), ".")) {
				dirs = append(dirs, file.Name())
			}
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
