package wiper

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Algorithm defines the interface for wiping algorithms
type Algorithm interface {
	Wipe(tempDir string, targetBytes int64, progressChan chan<- Progress, startTime time.Time) error
	NumPasses() int
}

// GetAlgorithm returns the appropriate algorithm for the given method
func GetAlgorithm(method WipeMethod) Algorithm {
	switch method {
	case SinglePassZeros:
		return &SinglePassAlgorithm{}
	case DoD522022M:
		return &DoDAlgorithm{}
	case Gutmann:
		return &GutmannAlgorithm{}
	default:
		return &SinglePassAlgorithm{}
	}
}

// SinglePassAlgorithm implements single-pass zero write
type SinglePassAlgorithm struct{}

func (a *SinglePassAlgorithm) NumPasses() int {
	return 1
}

func (a *SinglePassAlgorithm) Wipe(tempDir string, targetBytes int64, progressChan chan<- Progress, startTime time.Time) error {
	return writePass(tempDir, targetBytes, 0x00, 1, 1, "Single Pass (Zeros)", progressChan, startTime)
}

// DoDAlgorithm implements DoD 5220.22-M (3 passes)
type DoDAlgorithm struct{}

func (a *DoDAlgorithm) NumPasses() int {
	return 3
}

func (a *DoDAlgorithm) Wipe(tempDir string, targetBytes int64, progressChan chan<- Progress, startTime time.Time) error {
	// Pass 1: Write zeros
	err := writePass(tempDir, targetBytes, 0x00, 1, 3, "DoD Pass 1/3 (0x00)", progressChan, startTime)
	if err != nil {
		return err
	}

	// Pass 2: Write ones
	err = writePass(tempDir, targetBytes, 0xFF, 2, 3, "DoD Pass 2/3 (0xFF)", progressChan, startTime)
	if err != nil {
		return err
	}

	// Pass 3: Write random
	err = writePassRandom(tempDir, targetBytes, 3, 3, "DoD Pass 3/3 (Random)", progressChan, startTime)
	if err != nil {
		return err
	}

	return nil
}

// GutmannAlgorithm implements Gutmann method (35 passes)
type GutmannAlgorithm struct{}

func (a *GutmannAlgorithm) NumPasses() int {
	return 35
}

func (a *GutmannAlgorithm) Wipe(tempDir string, targetBytes int64, progressChan chan<- Progress, startTime time.Time) error {
	// Gutmann method: 4 random passes + 27 pattern passes + 4 random passes
	// For simplicity, we'll do: 4 random + 27 alternating patterns + 4 random

	// First 4 random passes
	for i := 1; i <= 4; i++ {
		err := writePassRandom(tempDir, targetBytes, i, 35, fmt.Sprintf("Gutmann Pass %d/35 (Random)", i), progressChan, startTime)
		if err != nil {
			return err
		}
	}

	// 27 pattern passes (we'll alternate between different bytes)
	patterns := []byte{0x55, 0xAA, 0x92, 0x49, 0x24, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}
	for i := 0; i < 27; i++ {
		pattern := patterns[i%len(patterns)]
		err := writePass(tempDir, targetBytes, pattern, i+5, 35, fmt.Sprintf("Gutmann Pass %d/35 (0x%02X)", i+5, pattern), progressChan, startTime)
		if err != nil {
			return err
		}
	}

	// Final 4 random passes
	for i := 32; i <= 35; i++ {
		err := writePassRandom(tempDir, targetBytes, i, 35, fmt.Sprintf("Gutmann Pass %d/35 (Random)", i), progressChan, startTime)
		if err != nil {
			return err
		}
	}

	return nil
}

// writePass writes a single pass with a specific byte pattern
func writePass(tempDir string, targetBytes int64, pattern byte, currentPass, totalPasses int, methodName string, progressChan chan<- Progress, startTime time.Time) error {
	const bufferSize = 1024 * 1024 // 1 MB buffer
	buffer := make([]byte, bufferSize)

	// Fill buffer with pattern
	for i := range buffer {
		buffer[i] = pattern
	}

	return writeBuffer(tempDir, targetBytes, buffer, currentPass, totalPasses, methodName, progressChan, startTime)
}

// writePassRandom writes a single pass with random data
func writePassRandom(tempDir string, targetBytes int64, currentPass, totalPasses int, methodName string, progressChan chan<- Progress, startTime time.Time) error {
	const bufferSize = 1024 * 1024 // 1 MB buffer
	buffer := make([]byte, bufferSize)

	return writeBuffer(tempDir, targetBytes, buffer, currentPass, totalPasses, methodName, progressChan, startTime)
}

// writeBuffer writes data to fill the target space
func writeBuffer(tempDir string, targetBytes int64, buffer []byte, currentPass, totalPasses int, methodName string, progressChan chan<- Progress, startTime time.Time) error {
	var totalWritten int64 = 0
	fileIndex := 0

	// For random passes, we need to regenerate the buffer each time
	isRandom := methodName == fmt.Sprintf("Gutmann Pass %d/35 (Random)", currentPass) ||
		methodName == "DoD Pass 3/3 (Random)"

	for totalWritten < targetBytes {
		// Create a new file for this chunk
		filePath := filepath.Join(tempDir, fmt.Sprintf("wipe_%d_pass%d.tmp", fileIndex, currentPass))
		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create wipe file: %w", err)
		}

		// Write up to target size
		remaining := targetBytes - totalWritten
		fileSize := int64(0)

		for fileSize < remaining {
			writeSize := int64(len(buffer))
			if fileSize+writeSize > remaining {
				writeSize = remaining - fileSize
			}

			// For random passes, generate new random data
			if isRandom {
				_, err = rand.Read(buffer[:writeSize])
				if err != nil {
					file.Close()
					return fmt.Errorf("failed to generate random data: %w", err)
				}
			}

			n, err := file.Write(buffer[:writeSize])
			if err != nil {
				file.Close()
				// Disk might be full, which is expected
				if err == io.ErrShortWrite || err.Error() == "no space left on device" {
					break
				}
				return fmt.Errorf("failed to write: %w", err)
			}

			fileSize += int64(n)
			totalWritten += int64(n)

			// Send progress update
			if progressChan != nil {
				elapsed := time.Since(startTime)
				progress := Progress{
					BytesWritten:  totalWritten,
					TotalBytes:    targetBytes * int64(totalPasses),
					CurrentPass:   currentPass,
					TotalPasses:   totalPasses,
					CurrentMethod: methodName,
					TimeElapsed:   elapsed,
				}

				// Estimate remaining time
				if totalWritten > 0 {
					bytesPerSecond := float64(totalWritten) / elapsed.Seconds()
					remainingBytes := (targetBytes * int64(totalPasses)) - totalWritten
					progress.EstimatedTime = time.Duration(float64(remainingBytes)/bytesPerSecond) * time.Second
				}

				progressChan <- progress
			}

			// Check if we've reached the target
			if totalWritten >= targetBytes {
				break
			}
		}

		file.Close()
		fileIndex++

		// If we couldn't write more, break
		if fileSize == 0 {
			break
		}
	}

	return nil
}
