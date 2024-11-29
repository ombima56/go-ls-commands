package listfiles

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"

	"go-ls-commands/colors"
)

// GetMaxFileSize calculates the maximum size among the files.
func GetMaxFileSize(files []os.FileInfo) int64 {
	var maxSize int64
	for _, file := range files {
		if file.Size() > maxSize {
			maxSize = file.Size()
		}
	}
	return maxSize
}

// // GetMaxFieldLength calculates the maximum length of a specified field (permissions, owner, group, size, modTime, fileName)
func GetMaxFieldLength(files []os.FileInfo, fieldType string) int {
	maxLength := 0

	for _, file := range files {
		var field string

		switch fieldType {
		case "permissions":
			field = FileModeToString(file.Mode())
		case "owner":
			stat := file.Sys().(*syscall.Stat_t)
			owner, _ := user.LookupId(strconv.Itoa(int(stat.Uid)))
			field = owner.Username
		case "group":
			stat := file.Sys().(*syscall.Stat_t)
			group, _ := user.LookupGroupId(strconv.Itoa(int(stat.Gid)))
			field = group.Name
		case "size":
			field = fmt.Sprintf("%d", file.Size())
		case "modTime":
			field = file.ModTime().Format("Jan _2 15:04")
		case "fileName":
			field = file.Name()
		}

		if len(field) > maxLength {
			maxLength = len(field)
		}
	}
	return maxLength
}

// PrintFileInfo prints the file info with dynamic field alignment and optional size omission.
func PrintFileInfo(path string, file os.FileInfo, maxSize int64, maxFieldLengths map[string]int, longFormat bool) {
	stat := file.Sys().(*syscall.Stat_t)
	owner, _ := user.LookupId(strconv.Itoa(int(stat.Uid)))
	group, _ := user.LookupGroupId(strconv.Itoa(int(stat.Gid)))

	// Get file mode (permissions), number of links, owner, group, size
	permissions := FileModeToString(file.Mode())
	numLinks := stat.Nlink
	modTime := file.ModTime().Format("Jan _2 15:04")
	color := colors.GetFileColor(file)

	// Use majMinSize to determine the size (or device information)
	sizeStr := majMinSize(stat, file)
	sizeStr = fmt.Sprintf("%-*s", maxFieldLengths["size"], sizeStr)

	// Prepare symlink info if the file is a symlink
	symlinkTarget := ""
	if file.Mode()&os.ModeSymlink != 0 && longFormat { // Only handle symlinks if long format is enabled
		fullPath := path
		if path != file.Name() {
			fullPath = path + "/" + file.Name() // Manually construct the full path
		}
		target, err := os.Readlink(fullPath)
		if err == nil {
			symlinkTarget = target
		} else {
			symlinkTarget = "<unresolved>"
		}
		// Adjust permissions display to indicate it's a symlink
		permissions = "l" + permissions[1:]
	}

	// Format the width dynamically based on the max length of fields
	permissionsStr := fmt.Sprintf("%-*s", maxFieldLengths["permissions"], permissions)
	linksStr := fmt.Sprintf("%*d", maxFieldLengths["links"], numLinks)
	ownerStr := fmt.Sprintf("%-*s", maxFieldLengths["owner"], owner.Username)
	groupStr := fmt.Sprintf("%-*s", maxFieldLengths["group"], group.Name)
	modTimeStr := fmt.Sprintf("%-*s", maxFieldLengths["modTime"], modTime)

	// Print file information in a single line, append symlink info if available
	if symlinkTarget != "" {
		// Print symlink info directly after the filename
		fmt.Printf("%s %s %s %s %s %s %s%s -> %s%s\n",
			permissionsStr,
			linksStr,
			ownerStr,
			groupStr,
			sizeStr,
			modTimeStr,
			color,
			file.Name(),
			symlinkTarget,
			colors.Reset,
		)
	} else {
		// Print regular file info
		fmt.Printf("%s %s %s %s %s %s %s%s%s\n",
			permissionsStr,
			linksStr,
			ownerStr,
			groupStr,
			sizeStr,
			modTimeStr,
			color,
			file.Name(),
			colors.Reset,
		)
	}
}

func majMinSize(stat *syscall.Stat_t, info os.FileInfo) string {
	size := stat.Rdev
	if info.Mode()&os.ModeDevice != 0 {
		major := uint64(size >> 8)
		minor := uint64(size & 0xff)
		return fmt.Sprintf("%d, %d", major, minor)
	}
	return fmt.Sprintf("%d", info.Size())
}
