package sort

import "github.com/ombima56/go-ls-commands/listfiles"

func SortReverse(files []listfiles.FileInfo) {
	fileLen := len(files)
	for i := 0; i < fileLen/2; i++ {
		files[i], files[fileLen-i-1] = files[fileLen-i-1], files[i]
	}
}