package cleaner

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mat/gowipeme/internal/platform/darwin"
)

// RecentFilesCleaner handles clearing recent files lists
type RecentFilesCleaner struct {
	recentDocsPath string
}

// NewRecentFilesCleaner creates a new recent files cleaner
func NewRecentFilesCleaner() *RecentFilesCleaner {
	rc := &RecentFilesCleaner{}

	// Get recent documents path
	path, err := darwin.GetRecentDocumentsPath()
	if err == nil {
		rc.recentDocsPath = path
	}

	return rc
}

// Name returns the name of this cleaner
func (rc *RecentFilesCleaner) Name() string {
	return "Recent Files"
}

// DryRun returns info about what will be cleared
func (rc *RecentFilesCleaner) DryRun() ([]string, error) {
	items := make([]string, 0)

	// Check if recent documents path exists
	if rc.recentDocsPath != "" {
		if _, err := os.Stat(rc.recentDocsPath); err == nil {
			// Count .sfl2 files (recent items)
			entries, err := os.ReadDir(rc.recentDocsPath)
			if err == nil {
				count := 0
				for _, entry := range entries {
					if !entry.IsDir() && filepath.Ext(entry.Name()) == ".sfl2" {
						count++
					}
				}
				if count > 0 {
					items = append(items, fmt.Sprintf("Application recent documents (%d apps)", count))
				}
			}
		}
	}

	// Also clear Finder recent items
	home, err := darwin.GetHomeDir()
	if err == nil {
		finderPlist := filepath.Join(home, "Library/Preferences/com.apple.finder.plist")
		if _, err := os.Stat(finderPlist); err == nil {
			items = append(items, "Finder recent items")
		}
	}

	// Recent servers (network locations)
	recentServers := filepath.Join(home, "Library/Application Support/com.apple.sharedfilelist/com.apple.LSSharedFileList.RecentServers.sfl2")
	if _, err := os.Stat(recentServers); err == nil {
		items = append(items, "Recent network servers")
	}

	// Recent hosts
	recentHosts := filepath.Join(home, "Library/Application Support/com.apple.sharedfilelist/com.apple.LSSharedFileList.RecentHosts.sfl2")
	if _, err := os.Stat(recentHosts); err == nil {
		items = append(items, "Recent hosts")
	}

	return items, nil
}

// Clean removes recent files lists
func (rc *RecentFilesCleaner) Clean() error {
	errors := make([]error, 0)

	// Clean application recent documents
	if rc.recentDocsPath != "" {
		if _, err := os.Stat(rc.recentDocsPath); err == nil {
			entries, err := os.ReadDir(rc.recentDocsPath)
			if err == nil {
				for _, entry := range entries {
					if !entry.IsDir() && filepath.Ext(entry.Name()) == ".sfl2" {
						filePath := filepath.Join(rc.recentDocsPath, entry.Name())
						err := os.Remove(filePath)
						if err != nil {
							errors = append(errors, err)
						}
					}
				}
			}
		}
	}

	// Clear recent servers
	home, err := darwin.GetHomeDir()
	if err == nil {
		recentServers := filepath.Join(home, "Library/Application Support/com.apple.sharedfilelist/com.apple.LSSharedFileList.RecentServers.sfl2")
		if _, err := os.Stat(recentServers); err == nil {
			_ = os.Remove(recentServers)
		}

		// Clear recent hosts
		recentHosts := filepath.Join(home, "Library/Application Support/com.apple.sharedfilelist/com.apple.LSSharedFileList.RecentHosts.sfl2")
		if _, err := os.Stat(recentHosts); err == nil {
			_ = os.Remove(recentHosts)
		}

		// Clear recent applications
		recentApps := filepath.Join(home, "Library/Application Support/com.apple.sharedfilelist/com.apple.LSSharedFileList.RecentApplications.sfl2")
		if _, err := os.Stat(recentApps); err == nil {
			_ = os.Remove(recentApps)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to clean some recent files: %v", errors)
	}

	return nil
}
