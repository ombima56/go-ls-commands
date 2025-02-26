package main

import (
	"fmt"
	"os"

	filepaths "go-ls-commands/filepath"
	"go-ls-commands/listfiles"
	"go-ls-commands/sorting"
)

func main() {
	args := os.Args[1:]
	var paths []string
	var flags []string

	// Separate paths and flags
	for _, arg := range args {
		if len(arg) > 0 && arg[0] == '-' && arg != "-" {
			flags = append(flags, arg)
		} else {
			// Expand tilde in paths
			expandedPath, err := filepaths.ExpandTilde(arg)
			if err != nil {
				fmt.Printf("Error expanding path %s: %v\n", arg, err)
				continue
			}
			paths = append(paths, expandedPath)
		}
	}

	// If no paths specified, use current directory
	if len(paths) == 0 {
		paths = append(paths, ".")
	}

	// Parse flags
	opts, err := listfiles.ValidateFlags(flags)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	sorting.SortFiles(paths)
	var validPaths []string

	// Process each path
	for _, path := range paths {
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

		if fileInfo.Mode()&os.ModeSymlink != 0 && !opts.LongFormat {
			//  a symlink
			target, err := listfiles.GetSymlinkTarget(path, fileInfo)
			if err != nil {
				fmt.Printf("path %q is not a symlink", path)
				continue
			}
			targetFileInfo, err := os.Lstat(target)
			if err == nil && targetFileInfo.IsDir() {
				fileInfo = targetFileInfo
			}
		}

		if fileInfo.IsDir() {
			validPaths = append(validPaths, path)
		} else {
			if opts.LongFormat {
				metadata := listfiles.NewFileMetadata()
				metadata.MaxSize = fileInfo.Size()
				listfiles.PrintFileInfo(path, fileInfo, metadata.MaxSize, metadata.MaxFieldLengths)
			} else {
				listfiles.PrintFileName(fileInfo)
				fmt.Println()
			}
		}
	}

	for i, path := range validPaths {
		// Print path header if we're listing multiple paths
		if len(paths) > 1 {
			if i > 0 {
				fmt.Println()
			}
			fmt.Printf("%s:\n", path)
		}

		// List directory contents
		listfiles.ListFiles(path, opts, i == 0)
	}
}
