package listfiles

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"

	"go-ls-commands/colors"
)

// FileModeToString converts file mode to string representation
func FileModeToString(mode os.FileMode) string {
	var perm string

	// Set the first character based on file type
	if mode&os.ModeCharDevice != 0 {
		perm = "c" // Character device
	} else if mode&os.ModeDevice != 0 {
		perm = "b" // Block device
	} else if mode.IsDir() {
		perm = "d"
	} else if mode&os.ModeSymlink != 0 {
		perm = "l"
	} else {
		perm = "-"
	}

	// Get the permission bits without the first character
	perms := mode.Perm().String()[1:]

	// Append the permission bits
	perm += perms

	// Check for special bits
	if mode&os.ModeSetuid != 0 {
		perm = perm[:3] + "s" + perm[4:]
	}
	if mode&os.ModeSetgid != 0 {
		perm = perm[:6] + "s" + perm[7:]
	}
	if mode&os.ModeSticky != 0 {
		perm = perm[:9] + "t"
	}

	return perm
}

// PrintFileName prints just the filename with appropriate color
func PrintFileName(file os.FileInfo) {
	color := colors.GetFileColor(file)
	fmt.Printf("%s%s%s ", color, file.Name(), colors.Reset)
}

// PrintFileInfo prints detailed file information
func PrintFileInfo(path string, file os.FileInfo, maxSize int64, maxFieldLengths map[string]int) {
	stat := file.Sys().(*syscall.Stat_t)
	owner, _ := user.LookupId(strconv.Itoa(int(stat.Uid)))
	group, _ := user.LookupGroupId(strconv.Itoa(int(stat.Gid)))

	// Get file attributes
	permissions := FileModeToString(file.Mode())
	numLinks := stat.Nlink
	modTime := file.ModTime().Format("Jan _2 15:04")
	color := colors.GetFileColor(file)

	// Construct full path for the file
	fullPath := path
	if path != file.Name() {
		fullPath = path + "/" + file.Name()
	}

	// Check for extended attributes
	extendedAttributes := ""
	if hasExtendedAttributes(fullPath) {
		extendedAttributes = "+"
	}

	// Format permissions with extended attributes included
	permWithExt := permissions + extendedAttributes

	// Pad permissions field to match max width
	if len(permWithExt) < maxFieldLengths["permissions"] {
		permWithExt = permWithExt + strings.Repeat(" ", maxFieldLengths["permissions"]-len(permWithExt))
	}

	// Format size or device info
	var sizeStr string
	if file.Mode()&os.ModeDevice != 0 || file.Mode()&os.ModeCharDevice != 0 {
		rdev := stat.Rdev

		// Extract major and minor numbers correctly using Linux conventions
		major := uint64((rdev>>8)&0xfff) | uint64((rdev>>32) & ^uint64(0xfff))
		minor := uint64(rdev&0xff) | uint64((rdev>>12) & ^uint64(0xff))

		// The standard ls uses fixed width with spacing between major and minor
		// The minor is right-aligned to the total field width
		majorStr := fmt.Sprintf("%d", major)
		minorStr := fmt.Sprintf("%d", minor)

		// Calculate field width for minor to match ls behavior
		// ls typically aligns with spaces between major and minor and pads to total width
		commaSpace := 2 // ", " takes 2 characters
		minorWidth := maxFieldLengths["size"] - len(majorStr) - commaSpace

		sizeStr = fmt.Sprintf("%s, %*s", majorStr, minorWidth, minorStr)
	} else {
		// For regular files - right align to the max field width
		sizeStr = fmt.Sprintf("%*d", maxFieldLengths["size"], file.Size())
	}

	// Handle symlink if needed
	symlinkTarget := getSymlinkTarget(path, file)

	// Format fields with proper alignment
	linksStr := fmt.Sprintf("%*d", maxFieldLengths["links"], numLinks)
	ownerStr := fmt.Sprintf("%-*s", maxFieldLengths["owner"], owner.Username)
	groupStr := fmt.Sprintf("%-*s", maxFieldLengths["group"], group.Name)
	modTimeStr := fmt.Sprintf("%-*s", maxFieldLengths["modTime"], modTime)

	// Print formatted output
	if symlinkTarget != "" {
		fmt.Printf("%s %s %s %s %s %s %s%s%s -> %s\n",
			permWithExt, linksStr, ownerStr, groupStr, sizeStr, modTimeStr,
			color, file.Name(), colors.Reset, symlinkTarget)
	} else {
		fmt.Printf("%s %s %s %s %s %s %s%s%s\n",
			permWithExt, linksStr, ownerStr, groupStr, sizeStr, modTimeStr,
			color, file.Name(), colors.Reset)
	}
}

// hasExtendedAttributes checks if the file has extended attributes
func hasExtendedAttributes(path string) bool {
	// Ensure the path is valid
	if path == "" {
		return false
	}

	buf := make([]byte, 0)
	size, err := syscall.Listxattr(path, buf)

	return err == nil && size > 0
}

// getSymlinkTarget gets the target of a symlink
func getSymlinkTarget(path string, file os.FileInfo) string {
	if file.Mode()&os.ModeSymlink == 0 {
		return "" // Not a symlink
	}

	fullPath := path
	if path != file.Name() {
		fullPath = path + "/" + file.Name()
	}

	target, err := os.Readlink(fullPath)
	if err != nil {
		return "<unresolved>"
	}

	return target
}

// GetSymlinkTarget gets the target of a symlink
func GetSymlinkTarget(path string, file os.FileInfo) (string, error) {
	if file.Mode()&os.ModeSymlink == 0 {
		return "", fmt.Errorf("not a symlink")
	}

	fullPath := path
	if path != file.Name() {
		fullPath = path + "/" + file.Name()
	}

	target, err := os.Readlink(fullPath)
	if err != nil {
		return "", fmt.Errorf("symlink target unresolved")
	}

	return target, nil
}

// updateFieldLengths updates the maximum field lengths map
func updateFieldLengths(path string, file os.FileInfo, maxLengths map[string]int) {
	// Get stat info for user/group lookups
	stat := file.Sys().(*syscall.Stat_t)

	// Check and update permissions length
	permissions := FileModeToString(file.Mode())

	// Account for potential '+' for extended attributes
	filePath := path
	if path != file.Name() {
		filePath = path + "/" + file.Name()
	}

	if hasExtendedAttributes(filePath) {
		if len(permissions)+1 > maxLengths["permissions"] {
			maxLengths["permissions"] = len(permissions) + 1
		}
	} else {
		if len(permissions) > maxLengths["permissions"] {
			maxLengths["permissions"] = len(permissions)
		}
	}

	// Check and update links length
	links := fmt.Sprintf("%d", stat.Nlink)
	if len(links) > maxLengths["links"] {
		maxLengths["links"] = len(links)
	}

	// Check and update owner length
	owner, _ := user.LookupId(strconv.Itoa(int(stat.Uid)))
	if len(owner.Username) > maxLengths["owner"] {
		maxLengths["owner"] = len(owner.Username)
	}

	// Check and update group length
	group, _ := user.LookupGroupId(strconv.Itoa(int(stat.Gid)))
	if len(group.Name) > maxLengths["group"] {
		maxLengths["group"] = len(group.Name)
	}

	// For device files, we need to simulate the exact output format that ls uses
	if file.Mode()&os.ModeDevice != 0 || file.Mode()&os.ModeCharDevice != 0 {
		rdev := stat.Rdev
		major := uint64((rdev>>8)&0xfff) | uint64((rdev>>32) & ^uint64(0xfff))
		minor := uint64(rdev&0xff) | uint64((rdev>>12) & ^uint64(0xff))

		// Calculate the space needed for "major, minor" format
		// Add 2 for ", " between major and minor
		deviceWidth := len(fmt.Sprintf("%d", major)) + 2 + len(fmt.Sprintf("%d", minor))

		if deviceWidth > maxLengths["size"] {
			maxLengths["size"] = deviceWidth
		}
	} else {
		// For regular files
		size := fmt.Sprintf("%d", file.Size())
		if len(size) > maxLengths["size"] {
			maxLengths["size"] = len(size)
		}
	}

	// Check and update modification time length
	modTime := file.ModTime().Format("Jan _2 15:04")
	if len(modTime) > maxLengths["modTime"] {
		maxLengths["modTime"] = len(modTime)
	}

	// Check and update filename length
	if len(file.Name()) > maxLengths["fileName"] {
		maxLengths["fileName"] = len(file.Name())
	}
}
