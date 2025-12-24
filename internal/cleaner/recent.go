package cleaner

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mat/gowipeme/internal/platform"
)

// RecentFilesCleaner handles clearing recent files lists
type RecentFilesCleaner struct {
	recentDocsPath string
}

// NewRecentFilesCleaner creates a new recent files cleaner
func NewRecentFilesCleaner() *RecentFilesCleaner {
	rc := &RecentFilesCleaner{}

	// Get recent documents path
	path, err := platform.GetRecentDocumentsPath()
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

	switch runtime.GOOS {
	case "darwin":
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
		home, err := platform.GetHomeDir()
		if err == nil {
			finderPlist := filepath.Join(home, "Library/Preferences/com.apple.finder.plist")
			if _, err := os.Stat(finderPlist); err == nil {
				items = append(items, "Finder recent items")
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
		}

	case "linux":
		// recently-used.xbel
		if rc.recentDocsPath != "" {
			if _, err := os.Stat(rc.recentDocsPath); err == nil {
				items = append(items, "Desktop recent items (recently-used.xbel)")
			}
			if _, err := os.Stat(rc.recentDocsPath + ".bak"); err == nil {
				items = append(items, "Desktop recent items backup (recently-used.xbel.bak)")
			}
		}

	case "windows":
		// %APPDATA%\\Microsoft\\Windows\\Recent (+ jump lists)
		if rc.recentDocsPath != "" {
			if info, err := os.Stat(rc.recentDocsPath); err == nil && info.IsDir() {
				if entries, err := os.ReadDir(rc.recentDocsPath); err == nil {
					items = append(items, fmt.Sprintf("Recent items (%d entries)", len(entries)))
				} else {
					items = append(items, "Recent items")
				}

				auto := filepath.Join(rc.recentDocsPath, "AutomaticDestinations")
				if _, err := os.Stat(auto); err == nil {
					items = append(items, "Jump Lists (AutomaticDestinations)")
				}
				custom := filepath.Join(rc.recentDocsPath, "CustomDestinations")
				if _, err := os.Stat(custom); err == nil {
					items = append(items, "Jump Lists (CustomDestinations)")
				}
			}
		}
	}

	return items, nil
}

// Clean removes recent files lists
func (rc *RecentFilesCleaner) Clean() error {
	errors := make([]error, 0)

	switch runtime.GOOS {
	case "darwin":
		// Clean application recent documents
		if rc.recentDocsPath != "" {
			if _, err := os.Stat(rc.recentDocsPath); err == nil {
				entries, err := os.ReadDir(rc.recentDocsPath)
				if err == nil {
					for _, entry := range entries {
						if !entry.IsDir() && filepath.Ext(entry.Name()) == ".sfl2" {
							filePath := filepath.Join(rc.recentDocsPath, entry.Name())
							if err := os.Remove(filePath); err != nil {
								errors = append(errors, err)
							}
						}
					}
				}
			}
		}

		// Clear recent servers/hosts/apps
		home, err := platform.GetHomeDir()
		if err == nil {
			recentServers := filepath.Join(home, "Library/Application Support/com.apple.sharedfilelist/com.apple.LSSharedFileList.RecentServers.sfl2")
			if _, err := os.Stat(recentServers); err == nil {
				_ = os.Remove(recentServers)
			}

			recentHosts := filepath.Join(home, "Library/Application Support/com.apple.sharedfilelist/com.apple.LSSharedFileList.RecentHosts.sfl2")
			if _, err := os.Stat(recentHosts); err == nil {
				_ = os.Remove(recentHosts)
			}

			recentApps := filepath.Join(home, "Library/Application Support/com.apple.sharedfilelist/com.apple.LSSharedFileList.RecentApplications.sfl2")
			if _, err := os.Stat(recentApps); err == nil {
				_ = os.Remove(recentApps)
			}
		}

	case "linux":
		// Remove the freedesktop recent files list
		if rc.recentDocsPath != "" {
			_ = os.Remove(rc.recentDocsPath)
			_ = os.Remove(rc.recentDocsPath + ".bak")
		}

	case "windows":
		// Best-effort: clear contents of Recent + Jump Lists
		clearDir := func(dir string) {
			entries, err := os.ReadDir(dir)
			if err != nil {
				errors = append(errors, err)
				return
			}
			for _, entry := range entries {
				_ = os.RemoveAll(filepath.Join(dir, entry.Name()))
			}
		}

		if rc.recentDocsPath != "" {
			if info, err := os.Stat(rc.recentDocsPath); err == nil && info.IsDir() {
				// Clear Jump Lists first
				auto := filepath.Join(rc.recentDocsPath, "AutomaticDestinations")
				if info, err := os.Stat(auto); err == nil && info.IsDir() {
					clearDir(auto)
				}
				custom := filepath.Join(rc.recentDocsPath, "CustomDestinations")
				if info, err := os.Stat(custom); err == nil && info.IsDir() {
					clearDir(custom)
				}
				// Clear Recent root
				clearDir(rc.recentDocsPath)
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to clean some recent files: %v", errors)
	}

	return nil
}
