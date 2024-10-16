package listfiles

import (
	"fmt"
	"os"
	"os/user"
	"sort"
	"strconv"
	"strings"
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
	permission := FileModeToString(file.Mode())
	numLinks := stat.Nlink
	fileSize := file.Size()

	modTime := file.ModTime().Format("okt 15 04:1")

	fmt.Printf("%s %d %s %s %6d %s %s\n",
		permission,
		numLinks,
		owner.Username,
		group.Name,
		fileSize,
		modTime,
		file.Name(),
	)
}

func ListFiles(dir string, showDetails, recursive, includeHidden, reverseOrder, sortByTime bool) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var files []FileInfo

	// Collecting all file information
	for _, entry := range entries {
		if !includeHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}

		files = append(files, FileInfo{
			Name:    info.Name(),
			IsDir:   info.IsDir(),
			ModTime: info.ModTime(),
		})
	}

	// Sort files based on the specified criteria
	if sortByTime {
		sort.Slice(files, func(i, j int) bool {
			if reverseOrder {
				return files[i].ModTime.After(files[j].ModTime)
			}
			return files[i].ModTime.Before(files[j].ModTime)
		})
	} else {
		sort.Slice(files, func(i, j int) bool {
			if reverseOrder {
				return files[i].Name > files[j].Name
			}
			return files[i].Name < files[j].Name
		})
	}

	// Display the files
	for _, file := range files {
		if showDetails {
			fileType := "FILE"
			if file.IsDir {
				fileType = "DIR"
			}
			fmt.Printf("%s\t%s\n", fileType, file.Name)
		} else {
			fmt.Println(file.Name)
		}
	}
	return nil
}
