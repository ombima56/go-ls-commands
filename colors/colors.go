package colors

import (
	"os"
	"strings"
)

func Colors() map[string]string {
	lsColors := os.Getenv("LS_COLORS")
	colorMap := make(map[string]string)

	if lsColors == "" {
		return colorMap // Return empty map if LS_COLORS is not set
	}

	pairs := strings.Split(lsColors, ":")
	for _, pair := range pairs {
		if strings.Contains(pair, "=") {
			parts := strings.Split(pair, "=")
			if len(parts) == 2 {
				colorMap[parts[0]] = parts[1]
			}
		}
	}

	return colorMap
}

var colorMap = Colors()

const (
	Reset = "\033[0m" // Reset color
)

// GetFileColor determines the color for a file based on its type using the colorMap.
func GetFileColor(file os.FileInfo) string {
	if file.IsDir() {
		if color, ok := colorMap["di"]; ok {
			return "\033[" + color + "m"
		}
		return "\033[34m" // Default to blue if not found
	}

	// Check for executable files
	if file.Mode().Perm()&0o111 != 0 {
		if color, ok := colorMap["ex"]; ok {
			return "\033[" + color + "m"
		}
		return "\033[32m" // Default to green if executable color is not found
	}

	// Check for symbolic links (symlinks)
	if file.Mode()&os.ModeSymlink != 0 {
		if color, ok := colorMap["ln"]; ok {
			return "\033[" + color + "m"
		}
		return "\033[36m" // Default to light cyan if not found
	}

	// Check for block devices
	if file.Mode()&os.ModeDevice != 0 && file.Mode()&os.ModeCharDevice == 0 {
		if color, ok := colorMap["bd"]; ok {
			return "\033[" + color + "m"
		}
		return "\033[33m" // Default to yellow for block devices
	}

	// Check for character devices
	if file.Mode()&os.ModeCharDevice != 0 {
		if color, ok := colorMap["cd"]; ok {
			return "\033[" + color + "m"
		}
		return "\033[33m" // Default to yellow for character devices
	}

	// Check for named pipes (e.g., FIFO files)
	if file.Mode()&os.ModeNamedPipe != 0 {
		if color, ok := colorMap["pi"]; ok {
			return "\033[" + color + "m"
		}
		return "\033[31m" // Default to red if not found
	}

	// Fallback to reset if no specific color is found
	return Reset
}
