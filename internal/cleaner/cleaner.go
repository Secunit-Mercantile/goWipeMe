package cleaner

import (
	"fmt"
	"sort"
	"strings"
)

// Cleaner defines the interface for all cleaning operations
type Cleaner interface {
	// Name returns the human-readable name of this cleaner
	Name() string

	// DryRun returns a list of items that would be deleted without actually deleting them
	DryRun() ([]string, error)

	// Clean performs the actual cleaning operation
	Clean() error
}

// CleanResult holds the result of a cleaning operation
type CleanResult struct {
	CleanerName string
	ItemsCleaned int
	BytesFreed int64
	Error error
}

// CleanerManager manages multiple cleaners
type CleanerManager struct {
	cleaners []Cleaner
}

// NewCleanerManager creates a new cleaner manager
func NewCleanerManager() *CleanerManager {
	return &CleanerManager{
		cleaners: make([]Cleaner, 0),
	}
}

// AddCleaner adds a cleaner to the manager
func (cm *CleanerManager) AddCleaner(cleaner Cleaner) {
	cm.cleaners = append(cm.cleaners, cleaner)
}

// GetCleaners returns all registered cleaners
func (cm *CleanerManager) GetCleaners() []Cleaner {
	return cm.cleaners
}

// DryRunAll runs dry-run on all cleaners and returns a summary
func (cm *CleanerManager) DryRunAll() (map[string][]string, error) {
	results := make(map[string][]string)

	for _, cleaner := range cm.cleaners {
		items, err := cleaner.DryRun()
		if err != nil {
			return nil, fmt.Errorf("%s dry-run failed: %w", cleaner.Name(), err)
		}
		if len(items) > 0 {
			results[cleaner.Name()] = items
		}
	}

	return results, nil
}

// CleanAll runs all cleaners and returns results
func (cm *CleanerManager) CleanAll() []CleanResult {
	results := make([]CleanResult, 0, len(cm.cleaners))

	for _, cleaner := range cm.cleaners {
		result := CleanResult{
			CleanerName: cleaner.Name(),
		}

		// Get items before cleaning to count them
		items, err := cleaner.DryRun()
		if err != nil {
			result.Error = err
			results = append(results, result)
			continue
		}

		result.ItemsCleaned = len(items)

		// Perform actual cleaning
		err = cleaner.Clean()
		if err != nil {
			result.Error = err
		}

		results = append(results, result)
	}

	return results
}

// Summary returns a formatted summary of dry-run results
func Summary(dryRunResults map[string][]string) string {
	var sb strings.Builder

	totalItems := 0
	for _, items := range dryRunResults {
		totalItems += len(items)
	}

	sb.WriteString(fmt.Sprintf("Total items to clean: %d\n\n", totalItems))

	names := make([]string, 0, len(dryRunResults))
	for cleanerName := range dryRunResults {
		names = append(names, cleanerName)
	}
	sort.Strings(names)

	for _, cleanerName := range names {
		items := dryRunResults[cleanerName]
		sb.WriteString(fmt.Sprintf("%s (%d items):\n", cleanerName, len(items)))
		for _, item := range items {
			sb.WriteString(fmt.Sprintf("  - %s\n", item))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
