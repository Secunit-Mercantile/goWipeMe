package platform

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetHomeDir returns the user's home directory.
func GetHomeDir() (string, error) {
	return os.UserHomeDir()
}

// ExpandPath expands a leading "~" to the user's home directory.
//
// Supported:
// - "~"
// - "~/something" (or "~\\something" on Windows)
//
// "~user" is not supported.
func ExpandPath(path string) (string, error) {
	if path == "" || path[0] != '~' {
		return path, nil
	}

	home, err := GetHomeDir()
	if err != nil {
		return "", err
	}

	if path == "~" {
		return home, nil
	}

	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~\\") {
		return filepath.Join(home, path[2:]), nil
	}

	return "", fmt.Errorf("unsupported ~ expansion: %q", path)
}


