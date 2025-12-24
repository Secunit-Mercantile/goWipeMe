# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- CONTRIBUTING.md with development guidelines
- CHANGELOG.md with automatic generation
- docs/ folder with architecture documentation
- docs/IMPROVEMENTS.md roadmap for planned enhancements

## [1.0.0] - 2024-12-24

### Added
- Initial release
- TUI interface using Bubble Tea
- GUI interface using Wails + Svelte 5
- Backup and restore for browser/shell history
- Multi-browser support: Safari, Chrome, Firefox, Edge, Brave, Arc
- Multi-shell support: Bash, Zsh, Fish
- Cache cleaning with whitelist protection
- Recent files cleaning (macOS, Linux, Windows)
- Clipboard clearing
- Secure disk wiping with three algorithms:
  - Single Pass (zeros)
  - DoD 5220.22-M (3 passes)
  - Gutmann (35 passes)
- Two-phase wipe safety to prevent OS crashes
- Cross-platform support: macOS (Intel/ARM), Linux (AMD64/ARM64/RISC-V), Windows
- CI/CD with GitHub Actions
- Multi-platform release builds

### Security
- Backup directory uses 0700 permissions
- Backup files use 0600 permissions
- Safety buffer prevents disk from filling completely during wipe

---

<!--
This changelog is automatically updated by the release workflow.
Manual entries should be added under [Unreleased].
-->
