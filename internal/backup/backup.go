package backup

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/mat/gowipeme/internal/platform/darwin"
)

// BackupInfo contains metadata about a backup
type BackupInfo struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Items     []string  `json:"items"`
	Size      int64     `json:"size"`
}

// BackupManager handles backup and restore operations
type BackupManager struct {
	backupDir string
}

// NewBackupManager creates a new backup manager
func NewBackupManager() (*BackupManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	backupDir := filepath.Join(homeDir, ".gowipeme", "backups")
	if err := os.MkdirAll(backupDir, 0700); err != nil {
		return nil, err
	}

	return &BackupManager{backupDir: backupDir}, nil
}

// backupItem represents a file to backup
type backupItem struct {
	Name       string
	SourcePath string
}

// getBackupItems returns all items that can be backed up
func (bm *BackupManager) getBackupItems() []backupItem {
	var items []backupItem

	// Browser histories
	browsers := []struct {
		name   string
		getter func() (string, error)
	}{
		{"Safari History", darwin.GetSafariHistoryPath},
		{"Chrome History", darwin.GetChromeHistoryPath},
		{"Edge History", darwin.GetEdgeHistoryPath},
		{"Brave History", darwin.GetBraveHistoryPath},
		{"Arc History", darwin.GetArcHistoryPath},
	}

	for _, browser := range browsers {
		if path, err := browser.getter(); err == nil {
			if _, err := os.Stat(path); err == nil {
				items = append(items, backupItem{Name: browser.name, SourcePath: path})
			}
		}
	}

	// Firefox uses profiles, find first profile with places.sqlite
	if profilesPath, err := darwin.GetFirefoxProfilesPath(); err == nil {
		entries, err := os.ReadDir(profilesPath)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					placesPath := filepath.Join(profilesPath, entry.Name(), "places.sqlite")
					if _, err := os.Stat(placesPath); err == nil {
						items = append(items, backupItem{Name: "Firefox History", SourcePath: placesPath})
						break
					}
				}
			}
		}
	}

	// Shell histories
	shells := []struct {
		name   string
		getter func() (string, error)
	}{
		{"Bash History", darwin.GetBashHistoryPath},
		{"Zsh History", darwin.GetZshHistoryPath},
		{"Fish History", darwin.GetFishHistoryPath},
	}

	for _, shell := range shells {
		if path, err := shell.getter(); err == nil {
			if _, err := os.Stat(path); err == nil {
				items = append(items, backupItem{Name: shell.name, SourcePath: path})
			}
		}
	}

	return items
}

// CreateBackup creates a new backup of all cleanable items
func (bm *BackupManager) CreateBackup() (*BackupInfo, error) {
	timestamp := time.Now()
	backupID := timestamp.Format("2006-01-02_15-04-05")
	backupPath := filepath.Join(bm.backupDir, backupID)

	if err := os.MkdirAll(backupPath, 0700); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	items := bm.getBackupItems()
	var backedUpItems []string
	var totalSize int64

	// Create a manifest to track original paths
	manifest := make(map[string]string)

	for _, item := range items {
		// Create a safe filename from the item name
		safeFilename := sanitizeFilename(item.Name)
		destPath := filepath.Join(backupPath, safeFilename)

		size, err := copyFile(item.SourcePath, destPath)
		if err != nil {
			// Log but continue with other items
			continue
		}

		manifest[safeFilename] = item.SourcePath
		backedUpItems = append(backedUpItems, item.Name)
		totalSize += size
	}

	if len(backedUpItems) == 0 {
		// Clean up empty backup directory
		os.RemoveAll(backupPath)
		return nil, fmt.Errorf("no items found to backup")
	}

	// Save manifest
	manifestPath := filepath.Join(backupPath, "manifest.json")
	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to create manifest: %w", err)
	}
	if err := os.WriteFile(manifestPath, manifestData, 0600); err != nil {
		return nil, fmt.Errorf("failed to save manifest: %w", err)
	}

	info := &BackupInfo{
		ID:        backupID,
		Timestamp: timestamp,
		Items:     backedUpItems,
		Size:      totalSize,
	}

	// Save backup info
	infoPath := filepath.Join(backupPath, "info.json")
	infoData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to create info file: %w", err)
	}
	if err := os.WriteFile(infoPath, infoData, 0600); err != nil {
		return nil, fmt.Errorf("failed to save info file: %w", err)
	}

	return info, nil
}

// ListBackups returns all available backups
func (bm *BackupManager) ListBackups() ([]BackupInfo, error) {
	entries, err := os.ReadDir(bm.backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []BackupInfo{}, nil
		}
		return nil, err
	}

	var backups []BackupInfo

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		infoPath := filepath.Join(bm.backupDir, entry.Name(), "info.json")
		data, err := os.ReadFile(infoPath)
		if err != nil {
			continue
		}

		var info BackupInfo
		if err := json.Unmarshal(data, &info); err != nil {
			continue
		}

		backups = append(backups, info)
	}

	// Sort by timestamp, newest first
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Timestamp.After(backups[j].Timestamp)
	})

	return backups, nil
}

// RestoreBackup restores a backup by ID
func (bm *BackupManager) RestoreBackup(backupID string) error {
	backupPath := filepath.Join(bm.backupDir, backupID)

	// Read manifest
	manifestPath := filepath.Join(backupPath, "manifest.json")
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("backup not found or corrupted: %w", err)
	}

	var manifest map[string]string
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return fmt.Errorf("invalid manifest: %w", err)
	}

	var restoreErrors []string

	for filename, originalPath := range manifest {
		sourcePath := filepath.Join(backupPath, filename)

		// Ensure the destination directory exists
		destDir := filepath.Dir(originalPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			restoreErrors = append(restoreErrors, fmt.Sprintf("%s: %v", filename, err))
			continue
		}

		if _, err := copyFile(sourcePath, originalPath); err != nil {
			restoreErrors = append(restoreErrors, fmt.Sprintf("%s: %v", filename, err))
			continue
		}
	}

	if len(restoreErrors) > 0 {
		return fmt.Errorf("some items failed to restore: %v", restoreErrors)
	}

	return nil
}

// DeleteBackup deletes a backup by ID
func (bm *BackupManager) DeleteBackup(backupID string) error {
	backupPath := filepath.Join(bm.backupDir, backupID)
	return os.RemoveAll(backupPath)
}

// GetBackup returns info for a specific backup
func (bm *BackupManager) GetBackup(backupID string) (*BackupInfo, error) {
	infoPath := filepath.Join(bm.backupDir, backupID, "info.json")
	data, err := os.ReadFile(infoPath)
	if err != nil {
		return nil, fmt.Errorf("backup not found: %w", err)
	}

	var info BackupInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("invalid backup info: %w", err)
	}

	return &info, nil
}

// PreviewBackup returns what would be backed up without actually doing it
func (bm *BackupManager) PreviewBackup() ([]string, error) {
	items := bm.getBackupItems()
	var names []string
	for _, item := range items {
		names = append(names, item.Name)
	}
	return names, nil
}

// Helper functions

func sanitizeFilename(name string) string {
	// Replace spaces and special chars with underscores
	result := ""
	for _, c := range name {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' {
			result += string(c)
		} else {
			result += "_"
		}
	}
	return result
}

func copyFile(src, dst string) (int64, error) {
	sourceFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destFile.Close()

	size, err := io.Copy(destFile, sourceFile)
	if err != nil {
		return 0, err
	}

	return size, destFile.Sync()
}
