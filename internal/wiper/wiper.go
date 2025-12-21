package wiper

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// WipeMethod represents the wiping algorithm to use
type WipeMethod int

const (
	// SinglePassZeros writes zeros in a single pass
	SinglePassZeros WipeMethod = iota
	// DoD522022M uses DoD 5220.22-M standard (3 passes)
	DoD522022M
	// Gutmann uses Gutmann method (35 passes - extreme paranoia)
	Gutmann
)

// String returns the human-readable name of the wipe method
func (w WipeMethod) String() string {
	switch w {
	case SinglePassZeros:
		return "Single Pass (Zeros)"
	case DoD522022M:
		return "DoD 5220.22-M (3-Pass)"
	case Gutmann:
		return "Gutmann Method (35-Pass)"
	default:
		return "Unknown"
	}
}

// Description returns a description of the wipe method
func (w WipeMethod) Description() string {
	switch w {
	case SinglePassZeros:
		return "Fast, sufficient for SSDs and modern drives"
	case DoD522022M:
		return "US DoD standard, 3 passes (0x00, 0xFF, random)"
	case Gutmann:
		return "Maximum security, 35 passes (very slow)"
	default:
		return ""
	}
}

// Progress represents the current progress of a wiping operation
type Progress struct {
	BytesWritten  int64
	TotalBytes    int64
	CurrentPass   int
	TotalPasses   int
	CurrentMethod string
	TimeElapsed   time.Duration
	EstimatedTime time.Duration
}

// Percentage returns the completion percentage (0-100)
func (p Progress) Percentage() float64 {
	if p.TotalBytes == 0 {
		return 0
	}
	return float64(p.BytesWritten) / float64(p.TotalBytes) * 100
}

// Wiper handles secure disk wiping operations
type Wiper struct {
	Method     WipeMethod
	VolumePath string
}

// NewWiper creates a new wiper for the specified volume and method
func NewWiper(volumePath string, method WipeMethod) (*Wiper, error) {
	// Validate volume path exists
	info, err := os.Stat(volumePath)
	if err != nil {
		return nil, fmt.Errorf("invalid volume path: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("volume path must be a directory")
	}

	return &Wiper{
		Method:     method,
		VolumePath: volumePath,
	}, nil
}

// GetFreeSpace returns the available free space on the volume in bytes
func (w *Wiper) GetFreeSpace() (int64, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(w.VolumePath, &stat)
	if err != nil {
		return 0, fmt.Errorf("failed to get volume stats: %w", err)
	}

	// Available blocks * block size
	freeSpace := int64(stat.Bavail) * int64(stat.Bsize)
	return freeSpace, nil
}

// GetVolumeInfo returns information about the volume
func (w *Wiper) GetVolumeInfo() (totalSpace, freeSpace int64, err error) {
	var stat syscall.Statfs_t
	err = syscall.Statfs(w.VolumePath, &stat)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get volume stats: %w", err)
	}

	totalSpace = int64(stat.Blocks) * int64(stat.Bsize)
	freeSpace = int64(stat.Bavail) * int64(stat.Bsize)
	return totalSpace, freeSpace, nil
}

// WipeFreeSpace wipes all free space on the volume using a safe two-pass strategy.
//
// Two-Phase Strategy (prevents OS crashes from full disk):
// Phase 1: Fill disk to 90% (or leave 1GB, whichever is larger) with wipe data
// Phase 2: Delete half of the wipe files, then wipe the freed space + original buffer
//
// This ensures the OS always has breathing room and won't crash from a full disk.
func (w *Wiper) WipeFreeSpace(progressChan chan<- Progress) error {
	// Get free space
	freeSpace, err := w.GetFreeSpace()
	if err != nil {
		return err
	}

	// Calculate safety buffer: 10% of free space or 1GB, whichever is larger
	safetyBuffer := freeSpace / 10
	const minSafetyBuffer = 1024 * 1024 * 1024 // 1 GB
	if safetyBuffer < minSafetyBuffer {
		safetyBuffer = minSafetyBuffer
	}

	// Ensure we have enough space
	if freeSpace <= safetyBuffer {
		return fmt.Errorf("insufficient free space (need at least %s)", FormatBytes(safetyBuffer))
	}

	// Create temporary directory for wipe files
	tempDir := filepath.Join(w.VolumePath, ".gowipeme_temp")
	err = os.MkdirAll(tempDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	algorithm := GetAlgorithm(w.Method)
	startTime := time.Now()

	// PHASE 1: Fill most of the disk, leaving safety buffer
	phase1Target := freeSpace - safetyBuffer
	err = algorithm.Wipe(tempDir, phase1Target, progressChan, startTime)
	if err != nil {
		return fmt.Errorf("phase 1 failed: %w", err)
	}

	// PHASE 2: Delete some wipe files to free up space, then wipe the reserved area
	// Get list of all wipe files
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read temp directory: %w", err)
	}

	// Calculate how many files to delete (half of the safety buffer worth)
	deleteTarget := safetyBuffer / 2
	var deletedSpace int64 = 0
	var filesToKeep []string

	// Delete files until we've freed enough space
	for _, entry := range entries {
		if deletedSpace >= deleteTarget {
			filesToKeep = append(filesToKeep, filepath.Join(tempDir, entry.Name()))
			continue
		}

		filePath := filepath.Join(tempDir, entry.Name())
		info, err := entry.Info()
		if err == nil {
			deletedSpace += info.Size()
			os.Remove(filePath)
		}
	}

	// Now wipe the space we freed + the original safety buffer
	phase2Target := deletedSpace + safetyBuffer
	err = algorithm.Wipe(tempDir, phase2Target, progressChan, startTime)
	if err != nil {
		return fmt.Errorf("phase 2 failed: %w", err)
	}

	// Cleanup is handled by defer os.RemoveAll(tempDir)
	return nil
}

// FormatBytes formats bytes into human-readable format
func FormatBytes(bytes int64) string {
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

// GetHomeDir returns the user's home directory
func GetHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home, nil
}
