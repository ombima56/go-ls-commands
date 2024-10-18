package listfiles

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	Reset  = "\033[0m"
	Blue   = "\033[34m" // Directory color
	Green  = "\033[32m" // File color
	Yellow = "\033[33m" // Executable color (optional)
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

func ListFiles(path string, longFormat bool, allFiles bool, recursive bool) {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	if longFormat {
		var totalSize int64
		for _, file := range files {
			fileInfo, _ := file.Info()

			// Skip hidden files if the -a flag is not set
			if !allFiles && strings.HasPrefix(fileInfo.Name(), ".") {
				continue
			}

			// Get the file's block size and accumulate
			stat := fileInfo.Sys().(*syscall.Stat_t)
			totalSize += int64(stat.Blocks)
		}
		fmt.Printf("total %d\n", totalSize)
	}

	for _, file := range files {
		fileInfo, _ := file.Info()

		// Skip hidden files if the -a flag is not set
		if !allFiles && strings.HasPrefix(fileInfo.Name(), ".") {
			continue
		}

		if longFormat {
			PrintFileInfo(fileInfo)
		} else {
			printFileName(fileInfo)
		}
		// If the file is a directory and the recursive flag is set, call ListFiles recursively
		if recursive && fileInfo.IsDir() {
			fmt.Printf("\n%s:\n", fileInfo.Name())
			ListFiles(path+"/"+fileInfo.Name(), longFormat, allFiles, recursive)
		}
	}

	if !longFormat {
		fmt.Println()
	}
}

func ValidateFlags(args []string) (bool, bool, bool, error) {
	var longFlag, allFlag, recursiveFlag bool
	validFlagProvided := false

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") {
			arg = strings.TrimPrefix(arg, "-")
			for _, char := range arg {
				switch char {
				case 'l':
					longFlag = true
					validFlagProvided = true
				case 'a':
					allFlag = true
					validFlagProvided = true
				case 'R':
					recursiveFlag = true
					validFlagProvided = true
				default:
					return false, false, false, fmt.Errorf("Invalid flag: -%s", string(char))
				}
			}
		} else if strings.HasPrefix(arg, "--") {
			arg = strings.TrimPrefix(arg, "--")
			switch arg {
			case "l":
				longFlag = true
				validFlagProvided = true
			case "a":
				allFlag = true
				validFlagProvided = true
			case "recursive":
				recursiveFlag = true
				validFlagProvided = true
			default:
				return false, false, false, fmt.Errorf("Invalid flag: --%s", arg)
			}
		} else {
			return false, false, false, fmt.Errorf("Invalid argument: %s", arg)
		}
	}

	if !validFlagProvided {
		return false, false, false, fmt.Errorf("No valid flag provided. Use -l, -a, or -R.")
	}

	return longFlag, allFlag, recursiveFlag, nil
}
