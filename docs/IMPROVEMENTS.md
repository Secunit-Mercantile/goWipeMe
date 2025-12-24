# Prioritized Improvements Roadmap

This document tracks planned improvements for goWipeMe, prioritized by impact and effort.

## Critical (High Impact, Should Do First)

- [ ] **Add version management**
  - Create `internal/version/version.go` with ldflags injection
  - Update Makefile to inject version at build time
  - Update CI workflows to inject version from git tags
  - Add `--version` flag to TUI CLI

- [ ] **Add unit tests for core packages**
  - `internal/cleaner/*` - Test all cleaner implementations
  - `internal/backup/backup.go` - Test backup/restore operations
  - `internal/wiper/methods.go` - Test algorithm correctness
  - `internal/platform/*` - Test path resolution per-OS

- [ ] **Add input validation for backup IDs**
  - Prevent path traversal attacks in `RestoreBackup()`
  - Prevent path traversal attacks in `DeleteBackup()`
  - Validate backup ID format matches expected pattern

## Important (Medium Impact)

- [ ] **Add golangci-lint to CI**
  - Add `.golangci.yml` configuration
  - Add lint job to CI workflow
  - Fix any existing lint issues

- [ ] **Fix silent failures with proper logging**
  - Add actual logging in `CreateBackup()` when items fail to copy
  - Consider adding a structured logger package

- [ ] **Add error context to RestoreBackup**
  - Return detailed per-file errors, not just summary
  - Allow partial restore success reporting

- [ ] **Sync wails.json version with git tags**
  - Update `productVersion` during release process
  - Consider using a version file as single source of truth

- [ ] **Add gosec security scanning to CI**
  - Add security scan job to CI workflow
  - Address any findings

## Nice to Have (Lower Priority)

- [ ] **Add test coverage reporting**
  - Integrate codecov or similar
  - Add coverage badges to README
  - Set coverage thresholds

- [ ] **Add backup integrity verification**
  - Store checksums in manifest.json
  - Verify checksums on restore
  - Report corrupted files

- [ ] **Improve sanitizeFilename()**
  - Handle unicode edge cases
  - Add length limits
  - Consider using a well-tested library

- [ ] **Refactor GetAlgorithm() to registry pattern**
  - Allow registering new algorithms without modifying switch
  - Improve Open/Closed principle compliance

- [ ] **Add structured logging package**
  - Replace ad-hoc error collection
  - Add log levels (debug, info, warn, error)
  - Support both TUI and file logging

---

*Last reviewed: 2024-12-24*
