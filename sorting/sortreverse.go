package sorting

import "io/fs"

func SortReverse(files []fs.FileInfo) {
	fileLen := len(files)
	for i := 0; i < fileLen/2; i++ {
		files[i], files[fileLen-i-1] = files[fileLen-i-1], files[i]
	}
}
