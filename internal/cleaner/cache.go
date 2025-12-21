package cleaner

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mat/gowipeme/internal/platform/darwin"
)

// CacheCleaner handles cleaning application caches
type CacheCleaner struct {
	cachePath string
	// Whitelist of cache directories to preserve (system-critical)
	whitelist map[string]bool
}

// NewCacheCleaner creates a new cache cleaner
func NewCacheCleaner() *CacheCleaner {
	cc := &CacheCleaner{
		whitelist: make(map[string]bool),
	}

	// Get cache path
	path, err := darwin.GetCachesPath()
	if err == nil {
		cc.cachePath = path
	}

	// Whitelist system-critical caches
	cc.whitelist["com.apple.bird"] = true           // iCloud
	cc.whitelist["com.apple.notificationcenter"] = true
	cc.whitelist["com.apple.LaunchServices"] = true
	cc.whitelist["com.apple.iconservices"] = true
	cc.whitelist["com.apple.nsurlsessiond"] = true

	return cc
}

// Name returns the name of this cleaner
func (cc *CacheCleaner) Name() string {
	return "Application Caches"
}

// DryRun returns a list of cache directories that will be cleaned
func (cc *CacheCleaner) DryRun() ([]string, error) {
	items := make([]string, 0)

	if cc.cachePath == "" {
		return items, nil
	}

	// Check if cache directory exists
	if _, err := os.Stat(cc.cachePath); os.IsNotExist(err) {
		return items, nil
	}

	// Read cache directory
	entries, err := os.ReadDir(cc.cachePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache directory: %w", err)
	}

	// Filter out whitelisted caches
	for _, entry := range entries {
		if entry.IsDir() {
			// Skip whitelisted directories
			if cc.whitelist[entry.Name()] {
				continue
			}

			// Get size estimate
			cachePath := filepath.Join(cc.cachePath, entry.Name())
			size := getDirSize(cachePath)
			sizeStr := formatSize(size)

			items = append(items, fmt.Sprintf("%s (%s)", entry.Name(), sizeStr))
		}
	}

	return items, nil
}

// Clean removes application cache directories
func (cc *CacheCleaner) Clean() error {
	if cc.cachePath == "" {
		return fmt.Errorf("cache path not found")
	}

	// Check if cache directory exists
	if _, err := os.Stat(cc.cachePath); os.IsNotExist(err) {
		return nil // Nothing to clean
	}

	// Read cache directory
	entries, err := os.ReadDir(cc.cachePath)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	errors := make([]error, 0)

	// Remove each cache directory (except whitelisted)
	for _, entry := range entries {
		if entry.IsDir() {
			// Skip whitelisted directories
			if cc.whitelist[entry.Name()] {
				continue
			}

			cachePath := filepath.Join(cc.cachePath, entry.Name())
			err := os.RemoveAll(cachePath)
			if err != nil {
				errors = append(errors, fmt.Errorf("%s: %w", entry.Name(), err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to clean some caches: %v", errors)
	}

	return nil
}

// getDirSize calculates the total size of a directory
func getDirSize(path string) int64 {
	var size int64

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	if err != nil {
		return 0
	}

	return size
}

// formatSize formats bytes into human-readable format
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
