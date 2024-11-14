package sorting

import (
	"io/fs"
	"strings"
)

// BubbleSortLowercaseFirst sorts the FileInfo slice prioritizing names that start with lowercase letters.
func BubbleSortLowercaseFirst(files []fs.FileInfo) {
	n := len(files)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			// Compare and prioritize lowercase-first
			if shouldSwap(files[j].Name(), files[j+1].Name()) {
				// Swap
				files[j], files[j+1] = files[j+1], files[j]
			}
		}
	}
}

// Helper function to determine if two names should be swapped
func shouldSwap(name1, name2 string) bool {
	// Check if the first character is lowercase
	firstLower := TrimNotAlpha(strings.ToLower(name1))
	secondLower := TrimNotAlpha(strings.ToLower(name2))
	if (firstLower == secondLower) && firstLower != "" && secondLower != "" {
		return name1 < name2
	}

	// If both are the same case, sort lexicographically
	return firstLower > secondLower
}

func TrimNotAlpha(str string) string {
	s := ""
	for _, ch := range str {
		if !((ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9')) {
			continue
		}
		s += string(ch)
	}
	return s
}
