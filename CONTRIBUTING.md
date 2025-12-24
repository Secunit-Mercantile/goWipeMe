# Contributing to goWipeMe

Thank you for your interest in contributing to goWipeMe! This document provides guidelines and instructions for contributing.

## Code of Conduct

Be respectful and constructive in all interactions. We welcome contributors of all experience levels.

## Getting Started

### Prerequisites

- **Go 1.25+** - [Download](https://go.dev/dl/)
- **Bun** - [Install](https://bun.sh/) (for frontend)
- **Wails v2.11.0** - `go install github.com/wailsapp/wails/v2/cmd/wails@v2.11.0` (for GUI)

### Setup

```bash
# Clone the repository
git clone https://github.com/mat/gowipeme.git
cd gowipeme

# Install frontend dependencies
make frontend-deps

# Run TUI in development mode
make dev-tui

# Run GUI in development mode (hot reload)
make dev-gui
```

### Project Structure

```
├── cmd/                    # Entry points
│   ├── gowipeme/          # TUI main
│   └── gowipeme-gui/      # GUI main (dev)
├── internal/              # Core packages
│   ├── backup/            # Backup/restore logic
│   ├── cleaner/           # Cleaning framework
│   ├── gui/               # Wails backend
│   ├── platform/          # OS-specific paths
│   ├── tui/               # Terminal UI
│   └── wiper/             # Secure wiping
├── frontend/              # Svelte GUI
└── docs/                  # Documentation
```

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

### 2. Make Changes

- Follow existing code patterns
- Add tests for new functionality
- Update documentation if needed

### 3. Test Your Changes

```bash
# Run all tests
make test

# Build both interfaces
make build-tui
make build-gui
```

### 4. Commit Your Changes

We use [Conventional Commits](https://www.conventionalcommits.org/) for automatic changelog generation:

```bash
# Format: <type>(<scope>): <description>

git commit -m "feat(cleaner): add Opera browser support"
git commit -m "fix(wiper): correct progress calculation for multi-pass"
git commit -m "docs: update installation instructions"
git commit -m "test(backup): add unit tests for restore"
```

**Types:**
| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation only |
| `test` | Adding/updating tests |
| `refactor` | Code change that neither fixes a bug nor adds a feature |
| `perf` | Performance improvement |
| `chore` | Build process or auxiliary tool changes |
| `ci` | CI configuration changes |

**Scopes:** `backup`, `cleaner`, `wiper`, `platform`, `tui`, `gui`, `frontend`

### 5. Submit a Pull Request

1. Push your branch to GitHub
2. Open a Pull Request against `main`
3. Fill out the PR template
4. Wait for CI checks to pass
5. Request review

## Adding New Features

### Adding a New Cleaner

1. Create a new file in `internal/cleaner/` (e.g., `myapp.go`)
2. Implement the `Cleaner` interface:

```go
type MyAppCleaner struct{}

func (c *MyAppCleaner) Name() string {
    return "MyApp History"
}

func (c *MyAppCleaner) DryRun() ([]string, error) {
    // Return list of items that would be cleaned
}

func (c *MyAppCleaner) Clean() error {
    // Perform actual cleaning
}
```

3. Register in `internal/tui/app.go` and `internal/gui/app.go`
4. Add tests in `internal/cleaner/myapp_test.go`

### Adding a New Wipe Algorithm

1. Implement the `Algorithm` interface in `internal/wiper/methods.go`:

```go
type MyAlgorithm struct{}

func (a *MyAlgorithm) NumPasses() int {
    return 5
}

func (a *MyAlgorithm) Wipe(tempDir string, targetBytes int64, progressChan chan<- Progress, startTime time.Time) error {
    // Implement wiping logic
}
```

2. Add to `WipeMethod` enum and `GetAlgorithm()` switch
3. Add tests

### Adding Platform Support

1. Create `internal/platform/paths_<os>.go` with build tag
2. Implement all path getter functions
3. Test on target platform

## Code Style

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions
- Keep functions focused and small
- Handle errors explicitly

## Testing

### Running Tests

```bash
# All tests
make test

# Specific package
go test -v ./internal/cleaner/...

# With coverage
go test -cover ./internal/...
```

### Writing Tests

- Place tests in `*_test.go` files next to the code
- Use table-driven tests where appropriate
- Mock filesystem operations for unit tests
- Test edge cases and error conditions

## Documentation

- Update `README.md` for user-facing changes
- Update `docs/ARCHITECTURE.md` for structural changes
- Add inline comments for complex logic
- Update `CHANGELOG.md` is automatic via CI

## Questions?

- Open an issue for bugs or feature requests
- Start a discussion for questions or ideas

Thank you for contributing!
