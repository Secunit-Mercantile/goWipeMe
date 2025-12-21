# Release Process

This document describes how to create a new release of goWipeMe.

## Prerequisites

- All changes committed and pushed to `main`
- Tests passing locally: `make test`
- Both TUI and GUI builds working locally

## Creating a Release

### 1. Choose Version Number

Use [Semantic Versioning](https://semver.org/):
- **MAJOR.MINOR.PATCH** (e.g., 1.2.3)
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes

### 2. Tag and Push

Use the helper script:

```bash
./scripts/tag-release.sh 1.0.0
```

Or manually:

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### 3. GitHub Actions Workflow

The push triggers `.github/workflows/release.yml` which:

1. ✅ **Runs tests** on all code
2. 🔨 **Builds TUI** for:
   - macOS (Intel & Apple Silicon)
   - Linux (AMD64 & ARM64)
   - Windows (AMD64)
3. 🖥️ **Builds GUI** for:
   - macOS (Universal, Intel, Apple Silicon)
   - Linux (AMD64 & ARM64)
   - Windows (AMD64)
4. 📦 **Compresses** all binaries
5. 🚀 **Creates GitHub Release** with all assets

### 4. Monitor Workflow

Check the Actions tab:
```
https://github.com/YOUR_USERNAME/gowipeme/actions
```

Build takes approximately 15-20 minutes.

### 5. Release Published

Once complete, the release appears at:
```
https://github.com/YOUR_USERNAME/gowipeme/releases
```

## Build Matrix

### TUI Binaries

| Platform | Architecture | File |
|----------|--------------|------|
| macOS | Intel (x86_64) | `gowipeme-tui-macos-intel.tar.gz` |
| macOS | Apple Silicon (ARM64) | `gowipeme-tui-macos-apple-silicon.tar.gz` |
| Linux | AMD64 | `gowipeme-tui-linux-amd64.tar.gz` |
| Linux | ARM64 | `gowipeme-tui-linux-arm64.tar.gz` |
| Windows | AMD64 | `gowipeme-tui-windows-amd64.zip` |

### GUI Binaries

| Platform | Architecture | File |
|----------|--------------|------|
| macOS | Universal (Intel + ARM64) | `gowipeme-gui-macos-universal.tar.gz` |
| macOS | Intel (x86_64) | `gowipeme-gui-macos-intel.tar.gz` |
| macOS | Apple Silicon (ARM64) | `gowipeme-gui-macos-apple-silicon.tar.gz` |
| Linux | AMD64 | `gowipeme-gui-linux-amd64.tar.gz` |
| Linux | ARM64 | `gowipeme-gui-linux-arm64.tar.gz` |
| Windows | AMD64 | `gowipeme-gui-windows-amd64.zip` |

## Optimizations Applied

All binaries are built with:
- `-ldflags="-s -w"` - Strip debug info (smaller size)
- `-trimpath` - Remove file system paths (reproducible builds)
- Compression (tar.gz or zip)

## Troubleshooting

### Build Fails

1. Check the Actions logs for specific errors
2. Test builds locally:
   ```bash
   make build-tui
   make build-gui
   make test
   ```

### Tag Already Exists

Delete the tag and recreate:
```bash
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0
./scripts/tag-release.sh 1.0.0
```

### Wrong Version Tagged

1. Delete the release on GitHub
2. Delete the tag (see above)
3. Create correct tag

## Post-Release

1. ✅ Test downloads from the release page
2. ✅ Update documentation if needed
3. ✅ Announce the release

## Local Testing Before Release

Test the release workflow components locally:

```bash
# Build all TUI variants
GOOS=darwin GOARCH=amd64 make build-tui-release
GOOS=darwin GOARCH=arm64 make build-tui-release
GOOS=linux GOARCH=amd64 make build-tui-release
GOOS=windows GOARCH=amd64 make build-tui-release

# Build GUI variants
make build-gui-release          # Apple Silicon
make build-gui-universal        # Universal binary
```
