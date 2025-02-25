package filepaths

import (
	"os/user"
	"strings"
)

func GetParentDir(path string) string {
	// Remove trailing slashes from the path
	if path != "/" {
		path = strings.TrimRight(path, "/")
	}
	if path == "." {
		return ".."
	}

	// Find the last index of the separator.
	lastSlashIndex := strings.LastIndex(path, "/")

	// If there is no Slash in the path return "." as it the current directory
	if lastSlashIndex == -1 {
		return "."
	}

	// Return the substring before the last slash, which is the parent directory
	parentDir := path[:lastSlashIndex]
	if parentDir == "" {
		return "/"
	}
	return parentDir
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

// ExpandTilde replaces the tilde (~) in a path with the current user's home directory
func ExpandTilde(path string) (string, error) {
	// If the path doesn't start with a tilde, return it unchanged
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	// Get the current user
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	// Replace the tilde with the home directory
	if path == "~" {
		return currentUser.HomeDir, nil
	}

	// Handle paths like "~/documents" or "~username/documents"
	if path[1] == '/' {
		// Path like "~/documents"
		return strings.Replace(path, "~", currentUser.HomeDir, 1), nil
	} else {
		// Path like "~username/documents"
		parts := strings.SplitN(path[1:], "/", 2)
		username := parts[0]
		
		// Get the specified user
		u, err := user.Lookup(username)
		if err != nil {
			return "", err
		}
		
		if len(parts) == 1 {
			return u.HomeDir, nil
		}
		return u.HomeDir + "/" + parts[1], nil
	}
}
