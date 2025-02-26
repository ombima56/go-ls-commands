package listfiles

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
	"syscall"

	filepaths "go-ls-commands/filepath"
	"go-ls-commands/sorting"
)

// ListFiles handles the listing of files and directories with various formats
func ListFiles(path string, opts Options, isFirst bool) {
	// Only show the directory header if recursive flag is set
	if isFirst && opts.Recursive {
		fmt.Println(path + ":")
	}

	// List current directory contents
	serveDir(path, opts)

	// Normalize path for recursive operations
	if opts.Recursive && len(path) > 1 {
		// if strings.HasSuffix(path, "/") {
		// 	path = strings.TrimSuffix(path, "/")
		// }
		path = strings.TrimSuffix(path, "/")
	}

	// Handle recursive directory traversal
	if opts.Recursive {
		processRecursive(path, opts)
	}
}

// serveDir handles listing files in a directory
func serveDir(dir string, opts Options) {
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
	if opts.AllFiles {
		curDirInfo, _ := os.Stat(dir)
		fileInfos = append(fileInfos, CustomFileInfo{curDirInfo, "."})

		parentDir := filepaths.GetParentDir(dir)
		parentDirInfo, _ := os.Stat(parentDir)
		fileInfos = append(fileInfos, CustomFileInfo{parentDirInfo, ".."})
	}

	// Filter and add other files
	for _, file := range files {
		if !opts.AllFiles && strings.HasPrefix(file.Name(), ".") {
			continue // skip hidden files if -a is not set
		}
		fileInfos = append(fileInfos, file)
	}

	// Sort files based on options
	sortFiles(fileInfos, opts)

	// Print files with correct format
	printDirectory(dir, fileInfos, opts)
}

// processRecursive handles recursive directory traversal
func processRecursive(path string, opts Options) {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("cannot read directory '%s': %v\n", path, err)
		return
	}

	fileInfos := make([]fs.FileInfo, 0)
	for _, v := range files {
		fileInfo, _ := v.Info()
		fileInfos = append(fileInfos, fileInfo)
	}

	// Sort files based on options
	sortFiles(fileInfos, opts)

	// Find subdirectories for recursion
	var dirs []string
	for _, file := range fileInfos {
		// Check if the file is a symlink
		if file.Mode()&os.ModeSymlink != 0 {
			// Resolve the symlink
			target, err := os.Readlink(path + "/" + file.Name())
			if err != nil {
				fmt.Printf("Error resolving symlink: %v\n", err)
				continue
			}
			// Check if the target is a directory
			targetInfo, err := os.Lstat(target)
			if err == nil && targetInfo.IsDir() {
				// Add the symlinked directory to the list of directories
				dirs = append(dirs, file.Name())
			}
		} else if file.IsDir() && (opts.AllFiles || !strings.HasPrefix(file.Name(), ".")) {
			dirs = append(dirs, file.Name())
		}
	}

	// Process each subdirectory recursively
	for _, dirName := range dirs {
		fullPath := path + "/" + dirName
		fmt.Printf("\n%s:\n", fullPath)
		ListFiles(fullPath, opts, false)
	}
}

// sortFiles applies sorting based on the provided options
func sortFiles(fileInfos []os.FileInfo, opts Options) {
	// Default sort by name
	sorting.BubbleSortLowercaseFirst(fileInfos)

	// Override with time sort if requested
	if opts.SortByTime {
		sorting.SortTime(fileInfos)
	}

	// Reverse the order if requested
	if opts.ReverseSort {
		sorting.SortReverse(fileInfos)
	}
}

// printDirectory prints a directory's contents in the proper format
func printDirectory(dir string, fileInfos []os.FileInfo, opts Options) {
	// Get file metadata for formatting
	metadata := CalculateFileMetadata(dir, fileInfos)

	// Print total blocks if using long format
	if opts.LongFormat {
		var totalBlocks int64
		for _, file := range fileInfos {
			if stat, ok := file.Sys().(*syscall.Stat_t); ok {
				totalBlocks += int64(stat.Blocks)
			}
		}
		fmt.Printf("total %d\n", totalBlocks/2)
	}

	// Print each file
	for _, file := range fileInfos {
		if opts.LongFormat {
			PrintFileInfo(dir, file, metadata.MaxSize, metadata.MaxFieldLengths)
		} else {
			PrintFileName(file)
		}
	}

	// Add newline if not in long format
	if !opts.LongFormat {
		fmt.Println()
	}
}

// calculateFileMetadata calculates metadata needed for formatting
func CalculateFileMetadata(dir string, fileInfos []os.FileInfo) FileMetadata {
	metadata := NewFileMetadata()

	for _, file := range fileInfos {
		// Update maximum size
		if file.Size() > metadata.MaxSize {
			metadata.MaxSize = file.Size()
		}

		// Update field lengths
		updateFieldLengths(dir, file, metadata.MaxFieldLengths)
	}

	return metadata
}
