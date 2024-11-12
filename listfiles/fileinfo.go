package listfiles

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"

	"go-ls-commands/colors"
)

func GetMaxFileSize(files []os.FileInfo) int64 {
	var maxSize int64
	for _, file := range files {
		if file.Size() > maxSize {
			maxSize = file.Size()
		}
	}
	return maxSize
}

func PrintFileInfo(path string, file os.FileInfo, maxSize int64) {
	stat := file.Sys().(*syscall.Stat_t)
	owner, _ := user.LookupId(strconv.Itoa(int(stat.Uid)))
	group, _ := user.LookupGroupId(strconv.Itoa(int(stat.Gid)))

	// Get file mode (permissions), number of links, owner, group, size
	permissions := FileModeToString(file.Mode())
	numLinks := stat.Nlink
	fileSize := file.Size()
	modTime := file.ModTime().Format("Jan _2 15:04")
	color := colors.GetFileColor(file)

	// Prepare symlink info if the file is a symlink
	symlinkInfo := ""
	if file.Mode()&os.ModeSymlink != 0 {
		fullPath := ""
		if path == file.Name() {
			fullPath = path
		} else {
			fullPath = path + "/" + file.Name() // Manually construct the full path
			// fmt.Println(path)
		}
		target, err := os.Readlink(fullPath)
		if err == nil {
			fileinfo, err1 := os.Lstat(target)
			if err1 != nil {
				symlinkInfo = fmt.Sprintf(" -> %s", target)
			} else {

				colorlink := colors.GetFileColor(fileinfo)
				symlinkInfo = fmt.Sprintf(" -> %s%s%s", colorlink, target, colors.Reset)
			}

		} else {
			symlinkInfo = fmt.Sprintf(" -> %v", path)
		}
		// Adjust permissions display to indicate it's a symlink
		permissions = "l" + permissions[1:]
	}

	// Get the width need for the file column.
	width := len(fmt.Sprintf("%d", maxSize))

	// Print information in ls -l format
	fmt.Printf("%s %d %s %s %*d %s %s%s%s%s\n",
		permissions,
		numLinks,
		owner.Username,
		group.Name,
		width,
		fileSize,
		modTime,
		color,
		file.Name(),
		colors.Reset,
		symlinkInfo,
	)
}
