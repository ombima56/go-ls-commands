package listfiles

import (
	"os"
)

// func FileModeToString(mode os.FileMode) string {
// 	var perm string
// 	if mode.IsDir() {
// 		perm = "d"
// 	} else if mode&os.ModeCharDevice != 0 {
// 		perm = ""
// 	} else {
// 		perm = "-"
// 	}

// 	perms := []rune(mode.String())
// 	for _, r := range perms[1:10] {
// 		perm += string(r)
// 	}

// 	return perm
// }

func FileModeToString(mode os.FileMode) string {
	var perm string
	// Check if it's a directory
	if mode.IsDir() {
		perm = "d"
	} else if mode&os.ModeCharDevice != 0 {
		perm = "" // For character devices
	} else if mode&os.ModeSymlink != 0 {
		perm = "l" // For symlinks
	} else {
		perm = "-" // Regular file
	}

	// Get the file's permission bits
	perms := []rune(mode.String())

	for i, r := range perms[1:10] {
		// We check for setuid, setgid, and sticky bits after the regular permissions
		if i == 2 && mode&os.ModeSetuid != 0 {
			// Setuid bit
			perm += "s"
		} else if i == 5 && mode&os.ModeSetgid != 0 {
			// Setgid bit
			perm += "s"
		} else if i == 8 && mode&os.ModeSticky != 0 {
			// Sticky bit
			perm += "t"
		} else {
			// Regular permission bits (rwx)
			perm += string(r)
		}
	}

	return perm
}
