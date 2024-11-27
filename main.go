package main

import (
	"fmt"
	"os"
	"strings"

	"go-ls-commands/listfiles"
	"go-ls-commands/sorting"
)

func main() {
	args := os.Args[1:]
	var paths []string
	var flags []string

	// Separate paths and flags
	for _, arg := range args {
		if len(arg) > 0 && arg[0] == '-' {
			flags = append(flags, arg)
		} else {
			paths = append(paths, arg)
		}
	}

	// If no paths specified, use current directory
	if len(paths) == 0 {
		paths = append(paths, ".")
	}

	// Parse flags
	longFlag, allFlag, recursiveFlag, timeFlag, reverseFlag := false, false, false, false, false
	if len(flags) > 0 {
		var err error
		longFlag, allFlag, recursiveFlag, timeFlag, reverseFlag, err = listfiles.ValidateFlags(flags)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}
	sorting.SortFiles(paths)
	newpaths := make([]string, 0)

	// make a function that is goinf to check  whether your path is a global package or not...

	// Process each path
	for _, path := range paths {
		if len(path) > 1 && !isSpecial(path) && path[0] == '/' {
			fmt.Printf("ls: cannot access '%s': No such file or directory\n", path)
			return
		}
		// if len(path) > 1 && path != "link" && string(path[len(path)-1]) == "/" {
		// 	path = strings.Trim(path, "/")
		// }

		// Check if path exists
		fileInfo, err := os.Lstat(path)

		if os.IsNotExist(err) {
			fmt.Printf("ls: cannot access '%s': No such file or directory\n", path)
			continue
		}
		if err != nil {
			fmt.Printf("Error accessing %s: %v\n", path, err)
			continue
		}

		// If it's a file, just print its info
		if !fileInfo.IsDir() {
			if longFlag {
				maxSize := listfiles.GetMaxFileSize([]os.FileInfo{fileInfo})
				listfiles.PrintFileInfo(path, fileInfo, maxSize)
			} else {
				listfiles.PrintFileName(fileInfo)
				fmt.Println()
			}
			continue
		} else {
			newpaths = append(newpaths, path)
		}
	}
	for i, path := range newpaths {
		// Print path header if we're listing multiple paths
		if len(paths) > 1 {
			if i > 0 {
				fmt.Println()
			}
			fmt.Printf("%s:\n", path)
		}

		// List directory contents
		listfiles.ListFiles(path, longFlag, allFlag, recursiveFlag, timeFlag, reverseFlag, i == 0)
	}
}

func isSpecial(path string) bool {
	spc := []string{"/usr", "/bin", "/dev"}
	for _, s := range spc {
		if strings.HasPrefix(path, s) {
			return true
		}
	}
	return false
}
