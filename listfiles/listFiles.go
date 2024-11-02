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

func printFileName(file os.FileInfo) {
	color := Reset
	if file.IsDir() {
		color = Blue
	}
	fmt.Printf("%s%s%s ", color, file.Name(), Reset)
}

func ListFiles(path string, longFormat bool, allFiles bool, recursive bool, isFirst bool) {
	// Only show the ".: " header if recursive flag is set
	if isFirst && recursive {
		fmt.Println(".:")
	}

	// List current directory contents
	serveDir(path, longFormat, allFiles, false)

	if recursive {
		files, err := os.ReadDir(path)
		if err != nil {
			fmt.Printf("cannot read directory '%s': %v\n", path, err)
			return
		}

		// Sort directories for consistent output
		var dirs []string
		for _, file := range files {
			if file.IsDir() && (allFiles || !strings.HasPrefix(file.Name(), ".")) {
				dirs = append(dirs, file.Name())
			}
		}
		sort.Strings(dirs)

		// Process each directory
		for _, dirName := range dirs {
			fullPath := filepath.Join(path, dirName)
			// Convert absolute path to relative path for display
			displayPath := filepath.Join(".", strings.TrimPrefix(fullPath, filepath.Dir(path)))
			fmt.Printf("\n%s:\n", displayPath)

			ListFiles(fullPath, longFormat, allFiles, recursive, false)
		}
	}
}

func serveDir(dir string, longFormat bool, allFiles bool, showDirName bool) {
	f, err := os.OpenFile(dir, os.O_RDONLY, 0o666)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer f.Close()

	// Read all files in directory
	files, err := f.Readdir(0)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Convert to a slice of FileInfo structs for consistent sorting
	var fileInfos []os.FileInfo
	if allFiles {
		// Add . and .. first
		if curDirInfo, err := os.Stat(dir); err == nil {
			fileInfos = append(fileInfos, curDirInfo) // Add .
		}
		if parentDirInfo, err := os.Stat(filepath.Dir(dir)); err == nil {
			fileInfos = append(fileInfos, parentDirInfo) // Add ..
		}
	}

	// Add all other files
	for _, file := range files {
		if !allFiles && strings.HasPrefix(file.Name(), ".") {
			continue // Skip hidden files if -a is not set
		}
		fileInfos = append(fileInfos, file)
	}

	// Sort files by name
	sort.Slice(fileInfos, func(i, j int) bool {
		// Special handling for . and ..
		nameI := fileInfos[i].Name()
		nameJ := fileInfos[j].Name()

		// Always put . and .. first
		if nameI == "." {
			return true
		}
		if nameJ == "." {
			return false
		}
		if nameI == ".." {
			return nameJ != "."
		}
		if nameJ == ".." {
			return nameI == "."
		}

		return nameI < nameJ
	})

	// Calculate total blocks if in long format
	if longFormat {
		var totalBlocks int64
		for _, file := range fileInfos {
			if stat, ok := file.Sys().(*syscall.Stat_t); ok {
				totalBlocks += int64(stat.Blocks)
			}
		}
		fmt.Printf("total %d\n", totalBlocks/2)
	}

	// Print files
	for i, file := range fileInfos {
		name := file.Name()
		// Special handling for . and ..
		if i == 0 && allFiles {
			name = "."
		} else if i == 1 && allFiles {
			name = ".."
		}

		if longFormat {
			PrintFileInfo(file)
		} else {
			color := Reset
			if file.IsDir() {
				color = Blue
			}
			fmt.Printf("%s%s%s  ", color, name, Reset)
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

func ValidateFlags(args []string) (bool, bool, bool, error) {
	var longFlag, allFlag, recursiveFlag bool

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
				default:
					return false, false, false, fmt.Errorf("invalid option --%s", flagStr)
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
					default:
						return false, false, false, fmt.Errorf("invalid option -- '%c'", flag)
					}
				}
			}
		}
	}

	return longFlag, allFlag, recursiveFlag, nil
}

func PrintFileName(file os.FileInfo) {
	color := Reset
	if file.IsDir() {
		color = Blue
	}
	fmt.Printf("%s%s%s  ", color, file.Name(), Reset)
}
