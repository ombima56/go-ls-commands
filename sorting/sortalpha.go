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
// func shouldSwap(name1, name2 string) bool {
// 	// Check if the first character is lowercase
// 	firstLower := TrimNotAlpha(strings.ToLower(name1))
// 	secondLower := TrimNotAlpha(strings.ToLower(name2))
// 	if (firstLower == secondLower) && firstLower != "" && secondLower != "" {
// 		return name1 < name2
// 	}

// 	// If both are the same case, sort lexicographically
// 	return firstLower > secondLower
// }

func shouldSwap(name1, name2 string) bool {
	// Handle empty strings by always placing them last
	if name1 == "" && name2 != "" {
		return true
	}
	if name2 == "" && name1 != "" {
		return false
	}
	firstRune1 := []rune(name1)[0]
	firstRune2 := []rune(name2)[0]

	name1IsAlpha := isAlpha(firstRune1)
	name2IsAlpha := isAlpha(firstRune2)
	name1IsDigit := isDigit(firstRune1)
	name2IsDigit := isDigit(firstRune2)

	// Alphabetic names should come before numeric names
	if name1IsAlpha && !name2IsAlpha {
		return false
	}
	if !name1IsAlpha && name2IsAlpha {
		return true
	}

	// Numeric names should come before special character names
	if name1IsDigit && !name2IsDigit {
		return false
	}
	if !name1IsDigit && name2IsDigit {
		return true
	}

	// Within alphabetic names, prioritize lowercase-first
	if name1IsAlpha && name2IsAlpha {
		firstLower := strings.ToLower(name1[:1]) == name1[:1]
		secondLower := strings.ToLower(name2[:1]) == name2[:1]
		if firstLower && !secondLower {
			return false
		}
		if !firstLower && secondLower {
			return true
		}
	}

	// Fall back to lexicographical comparison
	return name1 > name2
}
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
func isAlpha(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
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
