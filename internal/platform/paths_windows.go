//go:build windows
// +build windows

package platform

import (
	"fmt"
	"os"
	"path/filepath"
)

func localAppData() (string, error) {
	if v := os.Getenv("LOCALAPPDATA"); v != "" {
		return v, nil
	}
	home, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "AppData", "Local"), nil
}

func appData() (string, error) {
	if v := os.Getenv("APPDATA"); v != "" {
		return v, nil
	}
	home, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "AppData", "Roaming"), nil
}

// Browser history paths
func GetSafariHistoryPath() (string, error) {
	return "", fmt.Errorf("Safari is not available on Windows")
}

func GetChromeHistoryPath() (string, error) {
	lad, err := localAppData()
	if err != nil {
		return "", err
	}
	return filepath.Join(lad, "Google", "Chrome", "User Data", "Default", "History"), nil
}

func GetFirefoxProfilesPath() (string, error) {
	ad, err := appData()
	if err != nil {
		return "", err
	}
	return filepath.Join(ad, "Mozilla", "Firefox", "Profiles"), nil
}

func GetEdgeHistoryPath() (string, error) {
	lad, err := localAppData()
	if err != nil {
		return "", err
	}
	return filepath.Join(lad, "Microsoft", "Edge", "User Data", "Default", "History"), nil
}

func GetBraveHistoryPath() (string, error) {
	lad, err := localAppData()
	if err != nil {
		return "", err
	}
	return filepath.Join(lad, "BraveSoftware", "Brave-Browser", "User Data", "Default", "History"), nil
}

func GetArcHistoryPath() (string, error) {
	return "", fmt.Errorf("Arc is not supported on Windows")
}

// Shell history paths
func GetBashHistoryPath() (string, error) {
	// Useful for Git Bash / MSYS environments
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
	lad, err := localAppData()
	if err != nil {
		return "", err
	}
	// Best-effort: clear user temp directory
	return filepath.Join(lad, "Temp"), nil
}

// Recent files paths
func GetRecentDocumentsPath() (string, error) {
	ad, err := appData()
	if err != nil {
		return "", err
	}
	return filepath.Join(ad, "Microsoft", "Windows", "Recent"), nil
}


