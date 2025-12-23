//go:build !windows
// +build !windows

package wiper

import (
	"fmt"
	"syscall"
)

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

