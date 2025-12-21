.PHONY: build-tui build-gui build-gui-release build-gui-universal install clean test dev-tui dev-gui frontend-deps sign-and-run help

# Go build flags for optimized release
LDFLAGS := -s -w
GCFLAGS := -trimpath

# Build the TUI binary
build-tui:
	@echo "Building goWipeMe TUI..."
	@go build -o gowipeme cmd/gowipeme/main.go
	@echo "✓ Built: ./gowipeme"

# Build optimized TUI binary (smaller, stripped)
build-tui-release:
	@echo "Building goWipeMe TUI (release)..."
	@go build -ldflags="$(LDFLAGS)" -trimpath -o gowipeme cmd/gowipeme/main.go
	@echo "✓ Built: ./gowipeme (optimized)"

# Build the GUI binary (requires Wails setup)
build-gui:
	@echo "Building goWipeMe GUI..."
	@wails build -clean -o gowipeme-gui
	@echo "✓ Built: ./build/bin/goWipeMe.app"

# Build optimized GUI for Apple Silicon (smaller, stripped)
build-gui-release:
	@echo "Building goWipeMe GUI (release, Apple Silicon)..."
	@wails build -clean -platform darwin/arm64 -ldflags="$(LDFLAGS)" -trimpath -o gowipeme-gui
	@echo "✓ Built: ./build/bin/goWipeMe.app (optimized for Apple Silicon)"

# Build universal binary (Apple Silicon + Intel)
build-gui-universal:
	@echo "Building goWipeMe GUI (universal binary)..."
	@wails build -clean -platform darwin/universal -ldflags="$(LDFLAGS)" -trimpath -o gowipeme-gui
	@echo "✓ Built: ./build/bin/goWipeMe.app (universal binary)"

# Install binaries to /usr/local/bin
install: build-tui-release
	@echo "Installing binaries..."
	@sudo cp gowipeme /usr/local/bin/
	@echo "✓ Installed to /usr/local/bin"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f gowipeme gowipeme-gui
	@rm -rf build/ dist/ frontend/dist/
	@echo "✓ Clean complete"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run TUI in development mode
dev-tui: build-tui
	@./gowipeme

# Run Wails dev server (GUI)
dev-gui:
	@echo "Starting Wails dev server..."
	@wails dev

# Install frontend dependencies with Bun
frontend-deps:
	@echo "Installing frontend dependencies..."
	@cd frontend && bun install
	@echo "✓ Frontend dependencies installed"

# Sign and run the GUI app (required for macOS to run without quarantine)
sign-and-run:
	@echo "Signing goWipeMe.app..."
	@codesign --force --sign - "./build/bin/goWipeMe.app/Contents/MacOS/"*
	@codesign --force --deep --sign - "./build/bin/goWipeMe.app"
	@xattr -dr com.apple.quarantine "./build/bin/goWipeMe.app"
	@echo "✓ App signed successfully"
	@echo "Opening goWipeMe.app..."
	@open "./build/bin/goWipeMe.app"

# Show help
help:
	@echo "goWipeMe Makefile Commands:"
	@echo ""
	@echo "  Development:"
	@echo "    make dev-tui        - Build and run TUI"
	@echo "    make dev-gui        - Run Wails dev server"
	@echo ""
	@echo "  Build:"
	@echo "    make build-tui      - Build the TUI binary"
	@echo "    make build-gui      - Build the GUI binary"
	@echo ""
	@echo "  Release (optimized):"
	@echo "    make build-tui-release    - Build optimized TUI (stripped)"
	@echo "    make build-gui-release    - Build optimized GUI for Apple Silicon"
	@echo "    make build-gui-universal  - Build universal binary (ARM64 + x86_64)"
	@echo ""
	@echo "  Other:"
	@echo "    make sign-and-run   - Sign and run the GUI app (macOS)"
	@echo "    make install        - Install binaries to /usr/local/bin"
	@echo "    make clean          - Remove build artifacts"
	@echo "    make test           - Run tests"
	@echo "    make frontend-deps  - Install frontend dependencies"
	@echo ""

# Default target
.DEFAULT_GOAL := help
