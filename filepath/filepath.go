package filepaths

import (
	"strings"
)

func GetParentDir(path string) string {
	// Remove trailing slashes from the path
	path = strings.TrimRight(path, "/")

	// Find the last index of the separator.
	lastSlashIndex := strings.LastIndex(path, "/")

	// If there is no Slash in the path return "." as it the current directory
	if lastSlashIndex == -1 {
		return "."
	}

	// Return the substring before the last slash, which is the parent directory
	return path[:lastSlashIndex]
}

func JoinPaths(basePath string, additionalPaths ...string) string {
	// Ensurnig the base ends with a separator.
	if !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}

	// Join the additional paths
	for _, path := range additionalPaths {
		// Removind additional slashes from the paths
		path = strings.TrimLeft(path, "/")
		basePath += path + "/"
	}
	return strings.TrimRight(basePath, "/")
}
