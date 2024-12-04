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

func GetColumnWidths(files []os.FileInfo) (permissionsWidth, linksWidth, ownerWidth, groupWidth, sizeWidth, modTimeWidth int) {
	for _, file := range files {
		stat := file.Sys().(*syscall.Stat_t)
		owner, _ := user.LookupId(strconv.Itoa(int(stat.Uid)))
		group, _ := user.LookupGroupId(strconv.Itoa(int(stat.Gid)))

		permissions := FileModeToString(file.Mode())
		numLinks := len(fmt.Sprintf("%d", stat.Nlink))
		fileSize := len(majMinSize(stat, file))
		modTime := len(file.ModTime().Format("Jan _2 15:04"))

		// Update maximum widths dynamically
		if len(permissions) > permissionsWidth {
			permissionsWidth = len(permissions)
		}
		if numLinks > linksWidth {
			linksWidth = numLinks
		}
		if len(owner.Username) > ownerWidth {
			ownerWidth = len(owner.Username)
		}
		if len(group.Name) > groupWidth {
			groupWidth = len(group.Name)
		}
		if fileSize > sizeWidth {
			sizeWidth = fileSize
		}
		if modTime > modTimeWidth {
			modTimeWidth = modTime
		}
	}

	return
}

func PrintFileInfo(path string, file os.FileInfo, colWidths map[string]int) {
	stat := file.Sys().(*syscall.Stat_t)
	owner, _ := user.LookupId(strconv.Itoa(int(stat.Uid)))
	group, _ := user.LookupGroupId(strconv.Itoa(int(stat.Gid)))

	// Get file mode (permissions), number of links, owner, group, size
	permissions := FileModeToString(file.Mode())
	numLinks := stat.Nlink
	modTime := file.ModTime().Format("Jan _2 15:04")
	color := colors.GetFileColor(file)

	// Prepare symlink info if the file is a symlink
	symlinkInfo := ""
	if file.Mode()&os.ModeSymlink != 0 {
		fullPath := path + "/" + file.Name()
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

	fileSize := majMinSize(stat, file)

	// Dynamically adjust column widths based on provided widths
	fmt.Printf("%-*s %*d %-*s %-*s %*s %s %s%s%s%s\n",
		colWidths["permissions"], permissions,
		colWidths["links"], numLinks,
		colWidths["owner"], owner.Username,
		colWidths["group"], group.Name,
		colWidths["size"], fileSize,
		modTime,
		color,
		file.Name(),
		colors.Reset,
		symlinkInfo,
	)
}

func majMinSize(stat *syscall.Stat_t, info os.FileInfo) string {
	size := stat.Rdev
	if info.Mode()&os.ModeDevice != 0 {
		major := (uint64(size) >> 8) & 0xFF
		minor := (uint64(size) & 0xff) | ((uint64(size) >> 12) & 0xfff00)
		return fmt.Sprintf("%d,   %3d", major, minor)
	}
	return fmt.Sprintf("%d", info.Size())
}

func FileModeToString(mode os.FileMode) string {
	var fileType string

	switch {
	case mode.IsRegular():
		fileType = "-"
	case mode.IsDir():
		fileType = "d"
	case mode&os.ModeSymlink != 0:
		fileType = "l"
	case mode&os.ModeNamedPipe != 0:
		fileType = "p"
	case mode&os.ModeSocket != 0:
		fileType = "s"
	case mode&os.ModeDevice != 0:
		if mode&os.ModeCharDevice != 0 {
			fileType = "c"
		} else {
			fileType = "b"
		}
	default:
		fileType = "?"
	}

	permissions := []string{"---", "---", "---"}
	flags := []os.FileMode{0o400, 0o200, 0o100, 0o040, 0o020, 0o010, 0o004, 0o002, 0o001}

	for i, flag := range flags {
		if mode&flag != 0 {
			idx := i / 3
			perm := []string{"r", "w", "x"}[i%3]
			permissions[idx] = permissions[idx][:i%3] + perm + permissions[idx][i%3+1:]
		}
	}

	return fileType + permissions[0] + permissions[1] + permissions[2]
}
