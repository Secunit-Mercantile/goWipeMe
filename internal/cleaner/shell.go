package cleaner

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/atotto/clipboard"

	"github.com/mat/gowipeme/internal/platform"
)

// ShellCleaner handles cleaning shell history
type ShellCleaner struct {
	historyFiles map[string]string // shell name -> path
}

// NewShellCleaner creates a new shell history cleaner
func NewShellCleaner() *ShellCleaner {
	sc := &ShellCleaner{
		historyFiles: make(map[string]string),
	}

	// Discover shell history files
	sc.discoverHistoryFiles()

	return sc
}

// discoverHistoryFiles finds shell history files
func (sc *ShellCleaner) discoverHistoryFiles() {
	// Bash history
	if path, err := platform.GetBashHistoryPath(); err == nil {
		if _, err := os.Stat(path); err == nil {
			sc.historyFiles["Bash"] = path
		}
	}

	// Zsh history
	if path, err := platform.GetZshHistoryPath(); err == nil {
		if _, err := os.Stat(path); err == nil {
			sc.historyFiles["Zsh"] = path
		}
	}

	// Zsh sessions directory
	if path, err := platform.GetZshSessionsPath(); err == nil {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			sc.historyFiles["Zsh Sessions"] = path
		}
	}

	// Fish history
	if path, err := platform.GetFishHistoryPath(); err == nil {
		if _, err := os.Stat(path); err == nil {
			sc.historyFiles["Fish"] = path
		}
	}

	// PowerShell history (Windows)
	if runtime.GOOS == "windows" {
		// PSReadLine history (PowerShell 5+ / 7+)
		appData := os.Getenv("APPDATA")
		if appData == "" {
			if home, err := platform.GetHomeDir(); err == nil {
				appData = filepath.Join(home, "AppData", "Roaming")
			}
		}
		if appData != "" {
			psHistory := filepath.Join(appData, "Microsoft", "Windows", "PowerShell", "PSReadLine", "ConsoleHost_history.txt")
			if _, err := os.Stat(psHistory); err == nil {
				sc.historyFiles["PowerShell"] = psHistory
			}
		}
	}
}

// Name returns the name of this cleaner
func (sc *ShellCleaner) Name() string {
	return "Shell History"
}

// DryRun returns a list of shell history files that will be cleaned
func (sc *ShellCleaner) DryRun() ([]string, error) {
	items := make([]string, 0, len(sc.historyFiles))

	for shell, path := range sc.historyFiles {
		items = append(items, fmt.Sprintf("%s (%s)", shell, path))
	}

	return items, nil
}

// Clean removes shell history files
func (sc *ShellCleaner) Clean() error {
	errors := make([]error, 0)

	for shell, path := range sc.historyFiles {
		// Check if path exists
		info, err := os.Stat(path)
		if os.IsNotExist(err) {
			continue
		}

		// If it's a directory, remove all contents
		if info.IsDir() {
			err = os.RemoveAll(path)
			if err != nil {
				errors = append(errors, fmt.Errorf("%s: %w", shell, err))
				continue
			}
			// Recreate the directory
			err = os.MkdirAll(path, 0755)
			if err != nil {
				errors = append(errors, fmt.Errorf("%s: failed to recreate directory: %w", shell, err))
			}
		} else {
			// For files, truncate to zero length instead of deleting
			// This prevents shell warnings about missing history file
			err = os.Truncate(path, 0)
			if err != nil {
				errors = append(errors, fmt.Errorf("%s: %w", shell, err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to clean some shell histories: %v", errors)
	}

	return nil
}

// ClipboardCleaner handles clearing the system clipboard
type ClipboardCleaner struct{}

// NewClipboardCleaner creates a new clipboard cleaner
func NewClipboardCleaner() *ClipboardCleaner {
	return &ClipboardCleaner{}
}

// Name returns the name of this cleaner
func (cc *ClipboardCleaner) Name() string {
	return "Clipboard"
}

// DryRun returns info about clipboard cleaning
func (cc *ClipboardCleaner) DryRun() ([]string, error) {
	return []string{"System clipboard will be cleared"}, nil
}

// Clean clears the system clipboard using a cross-platform clipboard provider.
func (cc *ClipboardCleaner) Clean() error {
	if err := clipboard.WriteAll(""); err != nil {
		return fmt.Errorf("failed to clear clipboard: %w", err)
	}
	return nil
}
