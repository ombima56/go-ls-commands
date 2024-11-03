package logic

import (
	"github.com/ombima56/go-ls-commands/listfiles"
)

func SortTime(files []listfiles.FileInfo) {
	for i := 0; i < len(files); i++ {
		for j := 0; j < len(files)-i-1; j++ {
			if files[j].ModTime.Before(files[j+1].ModTime) {
				files[j], files[j+1] = files[j+1], files[j]
			}
		}
	}
}
