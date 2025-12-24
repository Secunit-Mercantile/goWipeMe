package cleaner

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mat/gowipeme/internal/platform"
)

// BrowserCleaner handles cleaning browser history
type BrowserCleaner struct {
	browsers map[string]string // browser name -> path
}

// NewBrowserCleaner creates a new browser cleaner
func NewBrowserCleaner() *BrowserCleaner {
	bc := &BrowserCleaner{
		browsers: make(map[string]string),
	}

	// Discover installed browsers
	bc.discoverBrowsers()

	return bc
}

// discoverBrowsers finds installed browsers and their history paths
func (bc *BrowserCleaner) discoverBrowsers() {
	// Safari
	if path, err := platform.GetSafariHistoryPath(); err == nil {
		if _, err := os.Stat(path); err == nil {
			bc.browsers["Safari"] = path
		}
	}

	// Chrome
	if path, err := platform.GetChromeHistoryPath(); err == nil {
		if _, err := os.Stat(path); err == nil {
			bc.browsers["Chrome"] = path
		}
	}

	// Chromium (Linux)
	if runtime.GOOS == "linux" {
		if path, err := platform.ExpandPath("~/.config/chromium/Default/History"); err == nil {
			if _, err := os.Stat(path); err == nil {
				bc.browsers["Chromium"] = path
			}
		}
	}

	// Firefox (check for profiles)
	if profilesPath, err := platform.GetFirefoxProfilesPath(); err == nil {
		profiles, err := filepath.Glob(filepath.Join(profilesPath, "*/places.sqlite"))
		if err == nil {
			for _, profile := range profiles {
				bc.browsers["Firefox"] = profile
				break // Use first profile for now
			}
		}
	}

	// Edge
	if path, err := platform.GetEdgeHistoryPath(); err == nil {
		if _, err := os.Stat(path); err == nil {
			bc.browsers["Edge"] = path
		}
	}

	// Brave
	if path, err := platform.GetBraveHistoryPath(); err == nil {
		if _, err := os.Stat(path); err == nil {
			bc.browsers["Brave"] = path
		}
	}

	// Arc
	if path, err := platform.GetArcHistoryPath(); err == nil {
		if _, err := os.Stat(path); err == nil {
			bc.browsers["Arc"] = path
		}
	}
}

// Name returns the name of this cleaner
func (bc *BrowserCleaner) Name() string {
	return "Browser History"
}

// DryRun returns a list of browsers that will be cleaned
func (bc *BrowserCleaner) DryRun() ([]string, error) {
	items := make([]string, 0, len(bc.browsers))

	for browser, path := range bc.browsers {
		items = append(items, fmt.Sprintf("%s (%s)", browser, path))
	}

	return items, nil
}

// Clean removes browser history files
func (bc *BrowserCleaner) Clean() error {
	errors := make([]error, 0)

	for browser, path := range bc.browsers {
		// Check if file still exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		// For SQLite databases, we can either delete the file or clear the tables
		// Deleting is simpler but may cause browser warnings on next launch
		// For now, we'll delete the file
		err := os.Remove(path)
		if err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", browser, err))
			continue
		}

		// Also remove associated files (WAL, SHM for SQLite)
		_ = os.Remove(path + "-wal")
		_ = os.Remove(path + "-shm")
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to clean some browsers: %v", errors)
	}

	return nil
}
