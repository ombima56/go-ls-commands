package listfiles

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
	"syscall"

	"go-ls-commands/sorting"
)

// ListFiles handles the listing of files and directories with various formats (long, recursive, etc.)
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

	// Read directory contents if recursive flag is set
	if recursive {
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

		// Sort files (lowercase first, time, and reverse sort based on flags)
		sorting.BubbleSortLowercaseFirst(fileInfos)
		if timeSort {
			sorting.SortTime(fileInfos)
		}
		if reverseSort {
			sorting.SortReverse(fileInfos)
		}

		// Calculate the maximum field lengths for longFormat output
		maxFieldLengths := calculateMaxFieldLengths(fileInfos)

		// Process each directory
		var dirs []string
		for _, file := range fileInfos {
			if file.IsDir() && (allFiles || !strings.HasPrefix(file.Name(), ".")) {
				dirs = append(dirs, file.Name())
			}
		}

		// Process each directory for recursion
		for _, dirName := range dirs {
			fullPath := path + "/" + dirName

			// Convert absolute path to relative path for display
			fmt.Printf("\n%s:\n", fullPath)

			// Recursively list files in the subdirectory
			ListFiles(fullPath, longFormat, allFiles, recursive, timeSort, reverseSort, false)
		}

		// Process each file (non-directory)
		for _, file := range fileInfos {
			if !file.IsDir() {
				// Print long format info or simple file name based on flag
				if longFormat {
					// Print file info with max field lengths
					maxSize := GetMaxFileSize(fileInfos)
					PrintFileInfo(path, file, maxSize, maxFieldLengths, longFormat)
				} else {
					// Print just the file name
					PrintFileName(file)
				}
			}
		}
	}
}

// calculateMaxFieldLengths calculates the maximum length for each file field (permissions, owner, etc.)
func calculateMaxFieldLengths(fileInfos []fs.FileInfo) map[string]int {
	maxFieldLengths := map[string]int{
		"permissions": 0,
		"links":       0,
		"owner":       0,
		"group":       0,
		"size":        0,
		"modTime":     0,
		"fileName":    0,
	}

	// Iterate over fileInfos to calculate max lengths for each field
	for _, file := range fileInfos {
		// Permissions length (e.g. "-rw-r--r--")
		permissions := file.Mode().String()
		if len(permissions) > maxFieldLengths["permissions"] {
			maxFieldLengths["permissions"] = len(permissions)
		}

		// Number of links (e.g. "1")
		links := fmt.Sprintf("%d", file.Mode())
		if len(links) > maxFieldLengths["links"] {
			maxFieldLengths["links"] = len(links)
		}

		// Owner and Group (Unix-specific, requires type assertion)
		if stat, ok := file.Sys().(*syscall.Stat_t); ok {
			owner := stat.Uid
			group := stat.Gid

			// For owner and group, you may need to lookup the UID/GID into usernames and groupnames,
			// which may require additional handling like using "user" or "group" packages
			// For now, we're just checking lengths as raw numbers (UID/GID).
			ownerStr := fmt.Sprintf("%d", owner)
			groupStr := fmt.Sprintf("%d", group)

			if len(ownerStr) > maxFieldLengths["owner"] {
				maxFieldLengths["owner"] = len(ownerStr)
			}
			if len(groupStr) > maxFieldLengths["group"] {
				maxFieldLengths["group"] = len(groupStr)
			}
		}

		// Size (e.g. "1234")
		size := fmt.Sprintf("%d", file.Size())
		if len(size) > maxFieldLengths["size"] {
			maxFieldLengths["size"] = len(size)
		}

		// Modification time (e.g. "Jan 2 15:04")
		modTime := file.ModTime().String()
		if len(modTime) > maxFieldLengths["modTime"] {
			maxFieldLengths["modTime"] = len(modTime)
		}

		// File name (e.g. "example.txt")
		name := file.Name()
		if len(name) > maxFieldLengths["fileName"] {
			maxFieldLengths["fileName"] = len(name)
		}
	}

	return maxFieldLengths
}
