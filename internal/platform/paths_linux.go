//go:build linux
// +build linux

package platform

import (
	"fmt"
	"os"
	"path/filepath"
)

// Browser history paths
func GetSafariHistoryPath() (string, error) {
	return "", fmt.Errorf("Safari is not available on Linux")
}

func GetChromeHistoryPath() (string, error) {
	// Google Chrome
	return ExpandPath("~/.config/google-chrome/Default/History")
}

func GetFirefoxProfilesPath() (string, error) {
	return ExpandPath("~/.mozilla/firefox")
}

func GetEdgeHistoryPath() (string, error) {
	return ExpandPath("~/.config/microsoft-edge/Default/History")
}

func GetBraveHistoryPath() (string, error) {
	return ExpandPath("~/.config/BraveSoftware/Brave-Browser/Default/History")
}

func GetArcHistoryPath() (string, error) {
	return "", fmt.Errorf("Arc is not supported on Linux")
}

// Shell history paths
func GetBashHistoryPath() (string, error) {
	return ExpandPath("~/.bash_history")
}

func GetZshHistoryPath() (string, error) {
	return ExpandPath("~/.zsh_history")
}

func GetZshSessionsPath() (string, error) {
	// Some zsh setups store session files here on Linux as well.
	return ExpandPath("~/.zsh_sessions")
}

func GetFishHistoryPath() (string, error) {
	return ExpandPath("~/.local/share/fish/fish_history")
}

// Application cache paths
func GetCachesPath() (string, error) {
	// Prefer XDG cache dir if set
	// (we use ExpandPath only for "~" expansion)
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return xdg, nil
	}
	return ExpandPath("~/.cache")
}

// Recent files paths
func GetRecentDocumentsPath() (string, error) {
	// Freedesktop recent files list
	if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
		return filepath.Join(xdg, "recently-used.xbel"), nil
	}
	return ExpandPath("~/.local/share/recently-used.xbel")
}


