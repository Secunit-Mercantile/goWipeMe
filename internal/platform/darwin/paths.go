package darwin

import (
	"os"
	"path/filepath"
)

// GetHomeDir returns the user's home directory
func GetHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home, nil
}

// ExpandPath expands ~ to the home directory
func ExpandPath(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	home, err := GetHomeDir()
	if err != nil {
		return "", err
	}

	if len(path) == 1 {
		return home, nil
	}

	return filepath.Join(home, path[1:]), nil
}

// Browser history paths
func GetSafariHistoryPath() (string, error) {
	return ExpandPath("~/Library/Safari/History.db")
}

func GetChromeHistoryPath() (string, error) {
	return ExpandPath("~/Library/Application Support/Google/Chrome/Default/History")
}

func GetFirefoxProfilesPath() (string, error) {
	return ExpandPath("~/Library/Application Support/Firefox/Profiles")
}

func GetEdgeHistoryPath() (string, error) {
	return ExpandPath("~/Library/Application Support/Microsoft Edge/Default/History")
}

func GetBraveHistoryPath() (string, error) {
	return ExpandPath("~/Library/Application Support/BraveSoftware/Brave-Browser/Default/History")
}

func GetArcHistoryPath() (string, error) {
	return ExpandPath("~/Library/Application Support/Arc/User Data/Default/History")
}

// Shell history paths
func GetBashHistoryPath() (string, error) {
	return ExpandPath("~/.bash_history")
}

func GetZshHistoryPath() (string, error) {
	return ExpandPath("~/.zsh_history")
}

func GetZshSessionsPath() (string, error) {
	return ExpandPath("~/.zsh_sessions")
}

func GetFishHistoryPath() (string, error) {
	return ExpandPath("~/.local/share/fish/fish_history")
}

// Application cache paths
func GetCachesPath() (string, error) {
	return ExpandPath("~/Library/Caches")
}

// Recent files paths
func GetRecentDocumentsPath() (string, error) {
	return ExpandPath("~/Library/Application Support/com.apple.sharedfilelist/com.apple.LSSharedFileList.ApplicationRecentDocuments")
}
