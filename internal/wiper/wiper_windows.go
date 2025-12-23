//go:build windows
// +build windows

package wiper

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	getDiskFreeSpaceEx = kernel32.NewProc("GetDiskFreeSpaceExW")
)

// GetFreeSpace returns the available free space on the volume in bytes
func (w *Wiper) GetFreeSpace() (int64, error) {
	var freeBytes int64
	var totalBytes int64
	var availableBytes int64

	volumePathPtr, err := syscall.UTF16PtrFromString(w.VolumePath)
	if err != nil {
		return 0, fmt.Errorf("failed to convert path: %w", err)
	}

	ret, _, _ := getDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(volumePathPtr)),
		uintptr(unsafe.Pointer(&freeBytes)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&availableBytes)),
	)

	if ret == 0 {
		return 0, fmt.Errorf("failed to get disk free space")
	}

	return freeBytes, nil
}

// GetVolumeInfo returns information about the volume
func (w *Wiper) GetVolumeInfo() (totalSpace, freeSpace int64, err error) {
	var freeBytes int64
	var totalBytes int64
	var availableBytes int64

	volumePathPtr, err := syscall.UTF16PtrFromString(w.VolumePath)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert path: %w", err)
	}

	ret, _, _ := getDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(volumePathPtr)),
		uintptr(unsafe.Pointer(&freeBytes)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&availableBytes)),
	)

	if ret == 0 {
		return 0, 0, fmt.Errorf("failed to get disk free space")
	}

	return totalBytes, freeBytes, nil
}

