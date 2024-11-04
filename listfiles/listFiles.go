package listfiles

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	Reset = "\033[0m"
	Blue  = "\033[34m" // Directory color
	Green = "\033[32m" // File color
)

type FileInfo struct {
	Name    string
	IsDir   bool
	ModTime time.Time
}

func FileModeToString(mode os.FileMode) string {
	var perm string
	if mode.IsDir() {
		perm = "d"
	} else {
		perm = "-"
	}

	perms := []rune(mode.String())
	for _, r := range perms[1:10] {
		perm += string(r)
	}

	return perm
}

func PrintFileInfo(file os.FileInfo) {
	stat := file.Sys().(*syscall.Stat_t)
	owner, _ := user.LookupId(strconv.Itoa(int(stat.Uid)))
	group, _ := user.LookupGroupId(strconv.Itoa(int(stat.Gid)))

	// Get file mode (permissions), number of links, owner, group, size
	permissions := FileModeToString(file.Mode())
	numLinks := stat.Nlink
	fileSize := file.Size()

	// Format modification time
	modTime := file.ModTime().Format("Jan 02 15:04")

	// Determine the output color based on whether it is a directory
	color := Reset
	if file.IsDir() {
		color = Blue
	}

	// Print information in ls -l format
	fmt.Printf("%s %d %s %s %6d %s %s%s%s\n",
		permissions,
		numLinks,
		owner.Username,
		group.Name,
		fileSize,
		modTime,
		color,
		file.Name(),
		Reset,
	)
}

func PrintFileName(file os.FileInfo) {
	color := Reset
	if file.IsDir() {
		color = Blue
	}
	fmt.Printf("%s%s%s ", color, file.Name(), Reset)
}

// func ListFiles(path string, longFormat bool, allFiles bool, recursive bool, timeSort bool, reverseSort bool, isFirst bool) {
// 	// Only show the ".: " header if recursive flag is set
// 	if isFirst && recursive {
// 		fmt.Println(".:")
// 	}

// 	// List current directory contents
// 	serveDir(path, longFormat, allFiles, timeSort, reverseSort, false)

// 	if recursive {
// 		files, err := os.ReadDir(path)
// 		if err != nil {
// 			fmt.Printf("cannot read directory '%s': %v\n", path, err)
// 			return
// 		}

// 		// Sort directories for consistent output
// 		var dirs []string
// 		for _, file := range files {
// 			if file.IsDir() && (allFiles || !strings.HasPrefix(file.Name(), ".")) {
// 				dirs = append(dirs, file.Name())
// 			}
// 		}
// 		// sort.Strings(dirs)
// 		// Sort the directories
// 		if reverseSort {
// 			sort.Sort(sort.Reverse(sort.StringSlice(dirs)))
// 		} else {
// 			sort.Strings(dirs)
// 		}

// 		// Process each directory
// 		for _, dirName := range dirs {
// 			fullPath := filepath.Join(path, dirName)
// 			// Convert absolute path to relative path for display
// 			displayPath := filepath.Join(".", strings.TrimPrefix(fullPath, filepath.Dir(path)))
// 			fmt.Printf("\n%s:\n", displayPath)

// 			ListFiles(fullPath, longFormat, allFiles, recursive, timeSort, reverseSort, false)
// 		}
// 	}
// }

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

		// Sort directories
		if reverseSort {
			sort.Sort(sort.Reverse(sort.StringSlice(dirs)))
		} else {
			sort.Strings(dirs)
		}

		// Process each directory
		for _, dirName := range dirs {
			fullPath := filepath.Join(path, dirName)
			// Convert absolute path to relative path for display
			displayPath := filepath.Join(".", dirName)
			fmt.Printf("\n%s:\n", displayPath)

			// Recursively list files in the subdirectory
			ListFiles(fullPath, longFormat, allFiles, recursive, timeSort, reverseSort, false)
		}
	}
}

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

// Helper function to create a custom FileInfo for . and ..
type customFileInfo struct {
	os.FileInfo
	name string
}

func (f customFileInfo) Name() string {
	return f.name
}

func ValidateFlags(args []string) (bool, bool, bool, bool, bool, error) {
	var longFlag, allFlag, recursiveFlag, timeFlag, reverseFlag bool

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			// Handle combined flags (-la)
			flagStr := strings.TrimPrefix(arg, "-")
			if strings.HasPrefix(flagStr, "-") {
				// Handle long flags (--long)
				flagStr = strings.TrimPrefix(flagStr, "-")
				switch flagStr {
				case "long":
					longFlag = true
				case "all":
					allFlag = true
				case "recursive":
					recursiveFlag = true
				case "time":
					timeFlag = true
				case "reverse":
					reverseFlag = true
				default:
					return false, false, false, false, false, fmt.Errorf("invalid option --%s", flagStr)
				}
			} else {
				// Handle short flags (-l)
				for _, flag := range flagStr {
					switch flag {
					case 'l':
						longFlag = true
					case 'a':
						allFlag = true
					case 'R':
						recursiveFlag = true
					case 't':
						timeFlag = true
					case 'r':
						reverseFlag = true
					default:
						return false, false, false, false, false, fmt.Errorf("invalid option -- '%c'", flag)
					}
				}
			}
		}
	}

	return longFlag, allFlag, recursiveFlag, timeFlag, reverseFlag, nil
}

// func PrintFileName(file os.FileInfo) {
// 	color := Reset
// 	if file.IsDir() {
// 		color = Blue
// 	}
// 	fmt.Printf("%s%s%s  ", color, file.Name(), Reset)
// }

