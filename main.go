package main

import (
	"fmt"
	"os"

	"go-ls-commands/listfiles"
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

	// Process each path
	for i, path := range paths {

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
				listfiles.PrintFileInfo(fileInfo)
			} else {
				listfiles.PrintFileName(fileInfo)
			}
			continue
		}

		// Print path header if we're listing multiple pathss
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
