package gui

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/mat/gowipeme/internal/backup"
	"github.com/mat/gowipeme/internal/cleaner"
	"github.com/mat/gowipeme/internal/wiper"
)

// App struct holds the application state
type App struct {
	ctx         context.Context
	cleanerMgr  *cleaner.CleanerManager
	wiper       *wiper.Wiper
	wiperMethod wiper.WipeMethod
	backupMgr   *backup.BackupManager
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// GetContext returns the application context
func (a *App) GetContext() context.Context {
	return a.ctx
}

// Startup is called when the app starts
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize cleaner manager
	a.cleanerMgr = cleaner.NewCleanerManager()
	a.cleanerMgr.AddCleaner(cleaner.NewBrowserCleaner())
	a.cleanerMgr.AddCleaner(cleaner.NewShellCleaner())
	a.cleanerMgr.AddCleaner(cleaner.NewCacheCleaner())
	a.cleanerMgr.AddCleaner(cleaner.NewRecentFilesCleaner())
	a.cleanerMgr.AddCleaner(cleaner.NewClipboardCleaner())

	// Initialize backup manager
	backupMgr, err := backup.NewBackupManager()
	if err == nil {
		a.backupMgr = backupMgr
	}
}

// CleanerInfo represents information about a cleaner
type CleanerInfo struct {
	Name  string   `json:"name"`
	Items []string `json:"items"`
	Count int      `json:"count"`
}

// GetCleanerStatus returns the current status of all cleaners
func (a *App) GetCleanerStatus() ([]CleanerInfo, error) {
	dryRunResults, err := a.cleanerMgr.DryRunAll()
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(dryRunResults))
	for cleanerName := range dryRunResults {
		names = append(names, cleanerName)
	}
	sort.Strings(names)

	infos := make([]CleanerInfo, 0, len(names))
	for _, cleanerName := range names {
		items := dryRunResults[cleanerName]
		infos = append(infos, CleanerInfo{
			Name:  cleanerName,
			Items: items,
			Count: len(items),
		})
	}

	return infos, nil
}

// RunCleaner runs all cleaners
func (a *App) RunCleaner() error {
	results := a.cleanerMgr.CleanAll()

	// Check for errors
	for _, result := range results {
		if result.Error != nil {
			return fmt.Errorf("%s: %w", result.CleanerName, result.Error)
		}
	}

	return nil
}

// WiperInfo represents information about the wiper
type WiperInfo struct {
	TotalSpace  int64  `json:"totalSpace"`
	FreeSpace   int64  `json:"freeSpace"`
	Volume      string `json:"volume"`
	Methods     []WipeMethodInfo `json:"methods"`
}

// WipeMethodInfo represents a wipe method
type WipeMethodInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GetWiperStatus returns wiper information
func (a *App) GetWiperStatus() (*WiperInfo, error) {
	// Get home directory as default volume
	homeDir, err := wiper.GetHomeDir()
	if err != nil {
		return nil, err
	}

	// Create temporary wiper to get volume info
	w, err := wiper.NewWiper(homeDir, wiper.SinglePassZeros)
	if err != nil {
		return nil, err
	}

	totalSpace, freeSpace, err := w.GetVolumeInfo()
	if err != nil {
		return nil, err
	}

	// Get available methods
	methods := []WipeMethodInfo{
		{
			ID:          int(wiper.SinglePassZeros),
			Name:        wiper.SinglePassZeros.String(),
			Description: wiper.SinglePassZeros.Description(),
		},
		{
			ID:          int(wiper.DoD522022M),
			Name:        wiper.DoD522022M.String(),
			Description: wiper.DoD522022M.Description(),
		},
		{
			ID:          int(wiper.Gutmann),
			Name:        wiper.Gutmann.String(),
			Description: wiper.Gutmann.Description(),
		},
	}

	return &WiperInfo{
		TotalSpace: totalSpace,
		FreeSpace:  freeSpace,
		Volume:     homeDir,
		Methods:    methods,
	}, nil
}

// RunWiper starts the wiping process
func (a *App) RunWiper(methodID int) error {
	// Get home directory
	homeDir, err := wiper.GetHomeDir()
	if err != nil {
		return err
	}

	// Create wiper with selected method
	method := wiper.WipeMethod(methodID)
	w, err := wiper.NewWiper(homeDir, method)
	if err != nil {
		return err
	}

	a.wiper = w
	a.wiperMethod = method

	// Run wiping (this will block)
	// In the real implementation, we'd want to run this in a goroutine
	// and send progress updates via runtime events
	progressChan := make(chan wiper.Progress, 10)
	go func() {
		for range progressChan {
			// In a real implementation, we'd emit runtime events here
			// runtime.EventsEmit(a.ctx, "wiper:progress", prog)
		}
	}()

	return w.WipeFreeSpace(progressChan)
}

// Greet returns a greeting message (example method)
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, welcome to goWipeMe!", name)
}

// BackupInfo represents information about a backup for the frontend
type BackupInfo struct {
	ID        string   `json:"id"`
	Timestamp string   `json:"timestamp"`
	Items     []string `json:"items"`
	Size      int64    `json:"size"`
}

// BackupPreview represents what will be backed up
type BackupPreview struct {
	Items []string `json:"items"`
}

// GetBackupPreview returns a preview of what will be backed up
func (a *App) GetBackupPreview() (*BackupPreview, error) {
	if a.backupMgr == nil {
		return nil, fmt.Errorf("backup manager not initialized")
	}

	items, err := a.backupMgr.PreviewBackup()
	if err != nil {
		return nil, err
	}

	return &BackupPreview{Items: items}, nil
}

// CreateBackup creates a new backup
func (a *App) CreateBackup() (*BackupInfo, error) {
	if a.backupMgr == nil {
		return nil, fmt.Errorf("backup manager not initialized")
	}

	info, err := a.backupMgr.CreateBackup()
	if err != nil {
		return nil, err
	}

	return &BackupInfo{
		ID:        info.ID,
		Timestamp: info.Timestamp.Format(time.RFC3339),
		Items:     info.Items,
		Size:      info.Size,
	}, nil
}

// ListBackups returns all available backups
func (a *App) ListBackups() ([]BackupInfo, error) {
	if a.backupMgr == nil {
		return nil, fmt.Errorf("backup manager not initialized")
	}

	backups, err := a.backupMgr.ListBackups()
	if err != nil {
		return nil, err
	}

	result := make([]BackupInfo, len(backups))
	for i, b := range backups {
		result[i] = BackupInfo{
			ID:        b.ID,
			Timestamp: b.Timestamp.Format(time.RFC3339),
			Items:     b.Items,
			Size:      b.Size,
		}
	}

	return result, nil
}

// RestoreBackup restores a backup by ID
func (a *App) RestoreBackup(backupID string) error {
	if a.backupMgr == nil {
		return fmt.Errorf("backup manager not initialized")
	}

	return a.backupMgr.RestoreBackup(backupID)
}

// DeleteBackup deletes a backup by ID
func (a *App) DeleteBackup(backupID string) error {
	if a.backupMgr == nil {
		return fmt.Errorf("backup manager not initialized")
	}

	return a.backupMgr.DeleteBackup(backupID)
}
