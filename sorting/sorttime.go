// In the sorting package
package sorting

import "io/fs"

// SortTime sorts files by modification time in descending order.
func SortTime(files []fs.FileInfo) {
	for i := 0; i < len(files); i++ {
		for j := 0; j < len(files)-i-1; j++ {
			if files[j].ModTime().Before(files[j+1].ModTime()) {
				files[j], files[j+1] = files[j+1], files[j]
			}
		}
	}
}
