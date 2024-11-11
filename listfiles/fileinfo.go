package listfiles

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func PrintFileInfo(file os.FileInfo) {
	stat := file.Sys().(*syscall.Stat_t)
	owner, _ := user.LookupId(strconv.Itoa(int(stat.Uid)))
	group, _ := user.LookupGroupId(strconv.Itoa(int(stat.Gid)))

	// Get file mode (permissions), number of links, owner, group, size
	permissions := FileModeToString(file.Mode())
	numLinks := stat.Nlink
	fileSize := file.Size()

	// Format modification time
	modTime := file.ModTime().Format("Jan _2 15:04")

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
