# Architecture Overview

This document describes the high-level architecture of goWipeMe.

## System Overview

goWipeMe is a cross-platform privacy tool that provides both Terminal (TUI) and Graphical (GUI) interfaces for:
- Backing up and restoring browser/shell history
- Clearing digital footprints (history, caches, recent files, clipboard)
- Secure wiping of disk free space

```
┌─────────────────────────────────────────────────────────────┐
│                      User Interfaces                         │
├──────────────────────────┬──────────────────────────────────┤
│      TUI (Bubble Tea)    │         GUI (Wails + Svelte)     │
│   cmd/gowipeme/main.go   │         main_gui.go              │
│   internal/tui/app.go    │   internal/gui/app.go            │
└──────────────────────────┴──────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Core Business Logic                     │
├─────────────────┬─────────────────┬─────────────────────────┤
│     Backup      │     Cleaner     │         Wiper           │
│ internal/backup │ internal/cleaner│    internal/wiper       │
└─────────────────┴─────────────────┴─────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Platform Abstraction                       │
│                    internal/platform                         │
├─────────────────┬─────────────────┬─────────────────────────┤
│     Darwin      │      Linux      │        Windows          │
│ paths_darwin.go │ paths_linux.go  │   paths_windows.go      │
└─────────────────┴─────────────────┴─────────────────────────┘
```

## Package Structure

### Entry Points

| Path | Description |
|------|-------------|
| `cmd/gowipeme/main.go` | TUI entry point |
| `cmd/gowipeme-gui/main.go` | GUI development entry point |
| `main_gui.go` | GUI production entry point (embeds frontend) |

### Core Packages

#### `internal/backup`
Manages backup and restore operations for browser and shell history.

**Key Types:**
- `BackupManager` - Handles backup directory and operations
- `BackupInfo` - Metadata about a backup (ID, timestamp, items, size)

**Key Features:**
- Creates timestamped backup directories in `~/.gowipeme/backups/`
- Uses manifest.json to track original file paths
- Supports preview (dry-run) before backup

#### `internal/cleaner`
Plugin-based cleaning framework using the Strategy pattern.

**Key Types:**
- `Cleaner` interface - Contract for all cleaners
- `CleanerManager` - Aggregates and runs multiple cleaners
- `CleanResult` - Result of a cleaning operation

**Implementations:**
| Cleaner | Responsibility |
|---------|---------------|
| `BrowserCleaner` | Safari, Chrome, Firefox, Edge, Brave, Arc history |
| `ShellCleaner` | Bash, Zsh, Fish history + clipboard |
| `CacheCleaner` | Application cache directories |
| `RecentFilesCleaner` | Recent file lists (OS-specific) |

#### `internal/wiper`
Secure disk wiping with multiple algorithms.

**Key Types:**
- `Wiper` - Main wiper struct
- `Algorithm` interface - Contract for wiping algorithms
- `Progress` - Progress reporting struct

**Algorithms:**
| Method | Passes | Description |
|--------|--------|-------------|
| SinglePassZeros | 1 | Fast, writes 0x00 |
| DoD522022M | 3 | DoD standard (0x00, 0xFF, random) |
| Gutmann | 35 | Maximum security |

**Safety Feature:** Two-phase wiping prevents OS crashes by maintaining 10% or 1GB buffer.

#### `internal/platform`
Cross-platform path resolution using build tags.

**Files:**
- `paths_common.go` - Shared utilities (ExpandPath, GetHomeDir)
- `paths_darwin.go` - macOS-specific paths
- `paths_linux.go` - Linux-specific paths
- `paths_windows.go` - Windows-specific paths

### UI Packages

#### `internal/tui`
Terminal UI using Bubble Tea framework.

**State Machine Views:**
```
menuView → cleanerView → resultsView
        → backupConfirmView → backupRunningView → resultsView
        → restoreSelectView → restoreConfirmView → restoreRunningView → resultsView
        → wiperMethodView → wiperConfirmView → wiperProgressView → resultsView
```

#### `internal/gui`
Wails backend exposing RPC methods for the Svelte frontend.

**Exported Methods:**
- `GetCleanerStatus()` / `RunCleaner()`
- `GetWiperStatus()` / `RunWiper(methodID)`
- `GetBackupPreview()` / `CreateBackup()` / `ListBackups()` / `RestoreBackup(id)`

### Frontend (`frontend/`)

Svelte 5 + Vite application with components:
- `Splash.svelte` - Loading screen
- `Home.svelte` - Feature selection
- `Backup.svelte` / `Restore.svelte` - Backup operations
- `Cleaner.svelte` - History cleaning
- `Wiper.svelte` - Secure disk wiping

## Design Patterns

### Strategy Pattern
Used in both `cleaner` and `wiper` packages to allow interchangeable algorithms.

### Manager Pattern
`CleanerManager` and `BackupManager` encapsulate complex operations.

### Dry-Run Pattern
All destructive operations support preview mode before execution.

### Progress Reporting
Long-running operations (wiper) use channels for async progress updates.

## Build System

See `Makefile` for available targets:
- `make dev-tui` / `make dev-gui` - Development
- `make build-tui` / `make build-gui` - Production builds
- `make test` - Run tests
- `make install` - Install to system

## Data Storage

| Location | Purpose |
|----------|---------|
| `~/.gowipeme/backups/` | Backup storage |
| `~/.gowipeme/backups/<id>/info.json` | Backup metadata |
| `~/.gowipeme/backups/<id>/manifest.json` | File path mappings |
