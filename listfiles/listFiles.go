package listfiles

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
	"time"
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

    // Print information in ls -l format
    fmt.Printf("%s %d %s %s %6d %s %s\n",
        permissions,
        numLinks,
        owner.Username,
        group.Name,
        fileSize,
        modTime,
        file.Name(),
    )
}

func printFileName(file os.FileInfo) {
	fmt.Printf("%s ", file.Name())
}

func ListFiles(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	var totalSize int64

	for _, file := range files {
		fileInfo, _ := file.Info()
		totalSize += fileInfo.Size()
	}

	fmt.Printf("total %d\n", totalSize/1024) // Display total size in kilobytes

	for _, file := range files {
		fileInfo, _ := file.Info()
		PrintFileInfo(fileInfo)
	}
}
