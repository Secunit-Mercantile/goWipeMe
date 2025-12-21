# goWipeMe - Project Context

This document captures the development history and technical details for continuing work on this project.

## Project Overview

**goWipeMe** is a cross-platform privacy tool for securely cleaning digital footprints and wiping free disk space. It provides both a TUI (Terminal User Interface) and GUI (Graphical User Interface).

**Repository:** `github.com/Secunit-Mercantile/goWipeMe`
**Module:** `github.com/mat/gowipeme`
**Current Version:** v1.0.4

## Tech Stack

| Component | Technology | Version |
|-----------|------------|---------|
| Language | Go | 1.21+ (developed on 1.25.5) |
| TUI Framework | Bubble Tea + Lipgloss | v1.3.10 / v1.1.0 |
| GUI Framework | Wails v2 | v2.11.0 |
| Frontend | Svelte 5 | 5.x |
| Frontend Build | Vite + Bun | Latest |
| Node.js | LTS | 22.x |

## Directory Structure

```
gowipeme/
├── cmd/
│   ├── gowipeme/              # TUI entry point
│   │   └── main.go
│   └── gowipeme-gui/          # GUI entry point (unused - see main_gui.go)
│       └── main.go
├── internal/
│   ├── backup/                # Backup & restore functionality
│   │   └── backup.go
│   ├── cleaner/               # History/cache cleaning
│   │   ├── cleaner.go         # Core interface & manager
│   │   ├── browser.go         # Browser history cleaning
│   │   ├── shell.go           # Shell history & clipboard
│   │   ├── cache.go           # Application caches
│   │   └── recent.go          # Recent files
│   ├── wiper/                 # Disk wiping
│   │   ├── wiper.go           # Core wiper with two-phase strategy
│   │   └── methods.go         # Wipe algorithms
│   ├── platform/
│   │   └── darwin/            # macOS-specific paths
│   │       └── paths.go
│   ├── tui/                   # Bubble Tea TUI
│   │   └── app.go
│   └── gui/                   # Wails backend
│       └── app.go
├── frontend/                  # Svelte 5 frontend
│   ├── src/
│   │   ├── App.svelte
│   │   ├── main.js
│   │   ├── assets/
│   │   │   └── gui-logo.png
│   │   └── components/
│   │       ├── Splash.svelte
│   │       ├── Home.svelte
│   │       ├── Backup.svelte
│   │       ├── Restore.svelte
│   │       ├── Cleaner.svelte
│   │       └── Wiper.svelte
│   └── wailsjs/               # Auto-generated Wails bindings
├── .github/workflows/
│   └── release.yml            # CI/CD for multi-platform builds
├── scripts/
│   └── tag-release.sh         # Helper for creating version tags
├── main_gui.go                # GUI main with embedded frontend assets
├── wails.json                 # Wails configuration
├── Makefile                   # Build automation
├── README.md                  # User documentation
├── RELEASE.md                 # Release process documentation
└── CONTEXT.md                 # This file
```

## Features Implemented

### 1. Backup & Restore
- **Location:** `internal/backup/backup.go`
- Creates timestamped backups in `~/.gowipeme/backups/`
- Backs up browser history databases and shell history files
- Supports listing, restoring, and deleting backups
- GUI components: `Backup.svelte`, `Restore.svelte`

### 2. Clear All History
- **Location:** `internal/cleaner/`
- **Browsers:** Safari, Chrome, Firefox, Edge, Brave, Arc
- **Shells:** Bash, Zsh, Fish
- **Other:** Application caches, recent files, clipboard
- Dry-run preview before deletion
- GUI component: `Cleaner.svelte`

### 3. Secure Wipe Free Space
- **Location:** `internal/wiper/`
- **Methods:**
  - Single Pass (Zeros) - Fast, sufficient for SSDs
  - DoD 5220.22-M (3-pass) - US DoD standard
  - Gutmann (35-pass) - Maximum paranoia mode
- **Two-Phase Safety Strategy:**
  - Phase 1: Fill to 90% (leave 10% or 1GB safety buffer)
  - Phase 2: Delete half wipe files, wipe freed space + original buffer
  - Prevents OS crashes from full disk
- GUI component: `Wiper.svelte`

### 4. GUI Features
- Splash screen with logo (2-second display, fade out)
- Home screen with 4 feature cards
- Progress tracking for all operations
- Backend in `internal/gui/app.go` exposes methods to frontend

## Key Implementation Details

### Svelte 5 Syntax (Critical)
The frontend uses Svelte 5 with the new runes API:
```javascript
// main.js - Use mount() not new App()
import { mount } from 'svelte'
const app = mount(App, { target: document.getElementById('app') })

// Components - Use $state() and $props()
let showSplash = $state(true)
let { onNavigate } = $props()

// Events - Use onclick not on:click
<button onclick={handleClick}>
```

### Build Tags for Testing
`main_gui.go` has a build tag to exclude from tests (embed directive fails without frontend/dist):
```go
//go:build !test
// +build !test
```

### macOS Code Signing
Run this after building to avoid quarantine issues:
```bash
make sign-and-run
# Or manually:
codesign --force --sign - "./build/bin/goWipeMe.app/Contents/MacOS/"*
codesign --force --deep --sign - "./build/bin/goWipeMe.app"
xattr -dr com.apple.quarantine "./build/bin/goWipeMe.app"
open "./build/bin/goWipeMe.app"
```

### Firefox History Path
Firefox uses profile directories. Iterate to find `places.sqlite`:
```go
profilesPath, _ := darwin.GetFirefoxProfilesPath()
entries, _ := os.ReadDir(profilesPath)
for _, entry := range entries {
    if entry.IsDir() {
        placesPath := filepath.Join(profilesPath, entry.Name(), "places.sqlite")
        // Check if exists...
    }
}
```

## CI/CD Pipeline

**Trigger:** Push tags matching `v*.*.*`

**Jobs:**
1. **Test:** `go test -v ./internal/...` (avoids main_gui.go embed issue)
2. **Build TUI:** 6 platforms (macOS Intel/ARM, Linux AMD64/ARM64/RISC-V, Windows)
3. **Build GUI macOS:** 3 variants (Universal, Intel, ARM)
4. **Build GUI Linux:** 3 variants (AMD64, ARM64, RISC-V)
5. **Build GUI Windows:** AMD64
6. **Create Release:** Combines all 14 artifacts with auto-generated notes

**Linux Dependency:** Uses `libwebkit2gtk-4.1-dev` (not 4.0)

## Build Commands

```bash
# Development
make dev-tui          # Build and run TUI
make dev-gui          # Run Wails dev server

# Production
make build-tui        # Standard TUI build
make build-tui-release # Optimized TUI (-ldflags="-s -w" -trimpath)
make build-gui        # Standard GUI build
make build-gui-release # Optimized GUI for Apple Silicon
make build-gui-universal # Universal binary (Intel + ARM)

# Testing
make test             # Run all tests

# macOS
make sign-and-run     # Sign and launch GUI
```

## Issues Fixed During Development

| Issue | Cause | Fix |
|-------|-------|-----|
| Blank GUI window | Svelte 5 installed but using Svelte 4 syntax | Updated to `mount()` and runes API |
| GitHub Actions test failure | `main_gui.go` embed directive without frontend/dist | Added `//go:build !test` tag |
| TUI build "not in std" | `cmd/gowipeme/main.go` path syntax | Changed to `./cmd/gowipeme` |
| GUI build "wails.json not found" | File was in `.gitignore` | Removed from `.gitignore`, committed file |
| Linux GUI build failure | Package `libwebkit2gtk-4.0-dev` not found | Updated to `libwebkit2gtk-4.1-dev` |
| "directory not found" in CI | `cmd/` directory not committed | Added `cmd/` to git |

## Future Enhancements

- [ ] Linux platform support (`internal/platform/linux/`)
- [ ] Windows platform support (`internal/platform/windows/`)
- [ ] Configuration file (YAML/TOML)
- [ ] SSD detection (recommend TRIM over multi-pass wiping)
- [ ] Smart browser detection (check if running before cleaning)
- [ ] Scheduled cleaning (cron integration)
- [ ] Custom cleaner plugins

## Commit History Summary

```
9ddb7f2 Add missing cmd/ directory with TUI and GUI entry points
d5fd7b4 Fix GitHub Actions build failures
d6c1a30 Fix GitHub Actions test failure
82bd7c3 Update README warning to use GitHub callout syntax
2fe2c9b Add comprehensive README with critical warnings
e816a82 Add complete goWipeMe application with GUI, TUI, and CI/CD
7a24e2b first commit
```

## Quick Start for New Session

1. **Check current state:**
   ```bash
   git status
   git log --oneline -5
   ```

2. **Build locally:**
   ```bash
   make build-tui && ./gowipeme        # Test TUI
   make build-gui && make sign-and-run  # Test GUI
   ```

3. **Create new release:**
   ```bash
   ./scripts/tag-release.sh X.Y.Z
   # Or manually:
   git tag -a vX.Y.Z -m "Release vX.Y.Z" && git push origin vX.Y.Z
   ```

4. **Monitor builds:**
   - Actions: https://github.com/Secunit-Mercantile/goWipeMe/actions
   - Releases: https://github.com/Secunit-Mercantile/goWipeMe/releases
