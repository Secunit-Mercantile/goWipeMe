# goWipeMe

A cross-platform TUI and GUI application for secure data wiping and privacy cleaning, written in Go.

## Features

### ✅ Current Features (Milestone 1 & 2)

#### 🧹 Clear All History
- **Browser History**: Safari, Chrome, Firefox, Edge, Brave, Arc
- **Shell History**: Bash, Zsh, Fish, and shell sessions
- **System Clipboard**: Clear clipboard contents
- **Dry-run preview**: See what will be deleted before confirming
- **Safe deletion**: Graceful handling of locked files and running applications

#### 💾 Secure Wipe Free Space
- **Multiple wipe methods**:
  - Single Pass (Zeros) - Fast, sufficient for SSDs
  - DoD 5220.22-M (3-Pass) - US DoD standard
  - Gutmann Method (35-Pass) - Maximum security paranoia mode
- **Real-time progress tracking**: Shows current pass, percentage, bytes written, and ETA
- **Volume information**: Displays total and free space before wiping
- **Safe operation**: Reserves 100MB to prevent disk filling completely

### 🚧 Coming Soon (Milestone 3-5)
- Application cache cleaning
- Recent files clearing
- GUI application with Wails + Svelte
- Cross-platform support (Linux, Windows)
- Scheduled cleaning
- Custom configuration

## Installation

### Requirements
- macOS 10.15+ (current version)
- Go 1.21+

### Build from Source

```bash
# Clone the repository
git clone https://github.com/mat/gowipeme
cd gowipeme

# Build the TUI
make build-tui

# Or run directly
make dev-tui
```

### Install System-Wide

```bash
make install
```

This installs the `gowipeme` binary to `/usr/local/bin`.

## Usage

### Terminal UI (TUI)

Run the interactive terminal interface:

```bash
./gowipeme
```

#### Navigation
- **Arrow keys** or **j/k**: Navigate menus
- **Enter**: Select option / Confirm action
- **q**: Go back / Quit (from main menu)
- **Ctrl+C**: Force quit (disabled during disk wiping)

### Clear All History Workflow
1. Select "Clear All History" from main menu
2. Review the dry-run preview showing what will be deleted
3. Press Enter to confirm or 'q' to cancel
4. View results summary

### Secure Wipe Free Space Workflow
1. Select "Secure Wipe Free Space" from main menu
2. Choose a wipe method:
   - Single Pass (Zeros) - Recommended for SSDs
   - DoD 5220.22-M - 3 passes for extra security
   - Gutmann Method - 35 passes (very slow!)
3. Review volume information and warnings
4. Press Enter to start wiping
5. Monitor real-time progress
6. View completion summary

## How It Works

### Browser History Cleaning
goWipeMe locates and removes browser history databases:
- Detects installed browsers automatically
- Safely deletes SQLite history files
- Removes associated WAL and SHM files
- Handles browser lock files gracefully

### Shell History Cleaning
Clears command history files for all common shells:
- Truncates history files to zero length (prevents shell warnings)
- Removes zsh session directories
- Preserves file structure for seamless shell operation

### Secure Free Space Wiping

#### Single Pass (Zeros)
Writes zeros (`0x00`) to all free space in a single pass. Fast and effective, especially for SSDs where data recovery is already difficult.

#### DoD 5220.22-M (3-Pass)
US Department of Defense standard:
1. **Pass 1**: Write zeros (`0x00`)
2. **Pass 2**: Write ones (`0xFF`)
3. **Pass 3**: Write random data

#### Gutmann Method (35-Pass)
Extreme security with 35 passes:
- 4 random passes
- 27 pattern passes (various byte patterns)
- 4 final random passes

**Note**: Modern drives with wear-leveling may not benefit from multiple passes. Single pass is usually sufficient.

## Project Structure

```
gowipeme/
├── cmd/
│   └── gowipeme/          # TUI entry point
├── internal/
│   ├── cleaner/           # History cleaning logic
│   │   ├── cleaner.go     # Core interface
│   │   ├── browser.go     # Browser history
│   │   └── shell.go       # Shell history & clipboard
│   ├── wiper/             # Disk wiping logic
│   │   ├── wiper.go       # Core wiper
│   │   └── methods.go     # Wipe algorithms
│   ├── platform/
│   │   └── darwin/        # macOS-specific code
│   └── tui/               # Bubble Tea UI
├── go.mod
├── Makefile
└── README.md
```

## Security Considerations

### ⚠️ Important Warnings
- **Irreversible**: All wiping and cleaning operations are permanent
- **Backups**: Always ensure you have backups before wiping
- **SSD Note**: Multiple-pass wiping is unnecessary for SSDs due to wear-leveling
- **Interruption**: Do not interrupt disk wiping once started
- **Permissions**: Some operations may require elevated privileges

### What goWipeMe Does NOT Do
- Does not require `sudo` for history cleaning
- Does not collect or transmit any data
- Does not modify system files or configurations
- Does not guarantee data recovery prevention on SSDs (use encryption instead)

## Development

### Makefile Commands

```bash
make build-tui      # Build TUI binary
make dev-tui        # Build and run TUI
make clean          # Remove build artifacts
make test           # Run tests (coming soon)
make install        # Install to /usr/local/bin
```

### Tech Stack
- **Language**: Go 1.21+
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **GUI Framework** (Coming): [Wails v2](https://wails.io/) + Svelte + Bun

## Roadmap

### ✅ Milestone 1: Core TUI & History Cleaning (Complete)
- [x] Project structure
- [x] Browser history cleaning
- [x] Shell history cleaning
- [x] Basic TUI with Bubble Tea
- [x] Dry-run preview

### ✅ Milestone 2: Disk Wiping (Complete)
- [x] Single-pass zero wipe
- [x] DoD 5220.22-M algorithm
- [x] Gutmann method
- [x] Progress reporting
- [x] TUI wiper views

### Milestone 3: Additional Cleaners
- [ ] Application cache cleaning
- [ ] Recent files clearing
- [ ] Advanced browser data (cookies, cache)

### Milestone 4: GUI Application
- [ ] Wails + Svelte frontend
- [ ] Shared backend logic
- [ ] macOS .app bundle packaging

### Milestone 5: Cross-Platform
- [ ] Linux support
- [ ] Windows support
- [ ] Platform-specific optimizations

## Contributing

This is currently a personal project, but contributions, issues, and feature requests are welcome!

## License

TBD

## Disclaimer

This software is provided "as is" without warranty of any kind. Use at your own risk. The authors are not responsible for any data loss or damage resulting from the use of this software.

**Always backup your data before using any wiping or cleaning tools.**
