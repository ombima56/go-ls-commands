package listfiles

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	filepaths "go-ls-commands/filepath"
	"go-ls-commands/sorting"
)

func serveDir(dir string, longFormat, allFiles, timeSort, reverseSort bool) {
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

	// Filter and process file list
	fileInfos := prepareFileList(dir, files, allFiles)

	// Sort files based on criteria
	sortFiles(fileInfos, timeSort, reverseSort)

	// Calculate total blocks if in long format
	var totalBlocks int64
	if longFormat {
		for _, file := range fileInfos {
			if stat, ok := file.Sys().(*syscall.Stat_t); ok {
				totalBlocks += int64(stat.Blocks)
			}
		}
		fmt.Printf("total %d\n", totalBlocks/2)
	}

	// Print file information
	printFiles(dir, fileInfos, longFormat)
}

// Helper function to prepare the file list
func prepareFileList(dir string, files []os.FileInfo, allFiles bool) []os.FileInfo {
	var fileInfos []os.FileInfo

	// Include . and .. if allFiles is set
	if allFiles {
		curDirInfo, _ := os.Stat(dir)
		fileInfos = append(fileInfos, customFileInfo{curDirInfo, "."})

		// parentDir := filepaths.GetParentDir(dir)
		parentDir := filepaths.GetParentDir(dir)
		parentDirInfo, _ := os.Stat(parentDir)
		fileInfos = append(fileInfos, customFileInfo{parentDirInfo, ".."})
	}

	// Add other files based on visibility
	for _, file := range files {
		if !allFiles && strings.HasPrefix(file.Name(), ".") {
			continue
		}
		fileInfos = append(fileInfos, file)
	}

	return fileInfos
}

// Helper function to sort files
func sortFiles(fileInfos []os.FileInfo, timeSort, reverseSort bool) {
	if timeSort {
		sorting.SortTime(fileInfos)
	} else {
		sorting.BubbleSortLowercaseFirst(fileInfos)
	}

	if reverseSort {
		sorting.SortReverse(fileInfos)
	}
}

// Helper function to print files
func printFiles(dir string, fileInfos []os.FileInfo, longFormat bool) {
	if longFormat {
		// Calculate column widths dynamically
		permissionsWidth, linksWidth, ownerWidth, groupWidth, sizeWidth, modTimeWidth := GetColumnWidths(fileInfos)

		// Create a map to pass column widths
		colWidths := map[string]int{
			"permissions": permissionsWidth,
			"links":       linksWidth,
			"owner":       ownerWidth,
			"group":       groupWidth,
			"size":        sizeWidth,
			"modTime":     modTimeWidth,
		}

		for _, file := range fileInfos {
			PrintFileInfo(dir, file, colWidths)
		}
	} else {
		for _, file := range fileInfos {
			PrintFileName(file)
		}
	}

	if !longFormat {
		fmt.Println()
	}
}
