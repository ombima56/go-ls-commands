package listfiles

import (
	"os"
)

func FileModeToString(mode os.FileMode) string {
	var perm string
	if mode.IsDir() {
		perm = "d"
	} else if mode&os.ModeCharDevice != 0 {
		perm = ""	
	} else {
		perm = "-"
	}

	perms := []rune(mode.String())
	for _, r := range perms[1:10] {
		perm += string(r)
	}

	return perm
}
