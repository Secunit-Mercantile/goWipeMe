.PHONY: build-tui build-gui build-gui-release build-gui-universal install clean test dev-tui dev-gui dev-gui-verbose dev-gui-status dev-gui-kill frontend-deps sign-and-run help

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

# Kill stuck dev server processes
dev-gui-kill:
	@echo "Killing dev server processes..."
	@pkill -f "wails dev" || true
	@pkill -f "vite" || true
	@pkill -f "bun.*dev" || true
	@echo "✓ Dev server processes killed"

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
	@echo "Checking for existing processes..."
	@if pgrep -f "wails dev" > /dev/null || pgrep -f "vite" > /dev/null || pgrep -f "bun.*dev" > /dev/null; then \
		echo "⚠️  Warning: Existing dev server processes detected!"; \
		echo "   Run 'make dev-gui-kill' first to clean up, or press Ctrl+C and kill them manually"; \
		echo ""; \
	fi
	@echo ""
	@echo "Timeline:"
	@echo "  1. 'Compiling frontend:' appears → Vite starts (usually <1 second)"
	@echo "  2. Vite shows 'ready' → Wails connects and builds (10-30 seconds)"
	@echo "  3. GUI window opens → You're ready!"
	@echo ""
	@echo "⚠️  After Vite says 'ready', WAIT 10-30 seconds for the GUI window to appear"
	@echo "   The process is working - Wails just doesn't show progress during connection"
	@echo ""
	@echo "📋 Check the terminal output above for any errors if window doesn't appear"
	@echo ""
	@wails dev

# Run Wails dev server with verbose output (for debugging)
dev-gui-verbose:
	@echo "Starting Wails dev server with verbose output..."
	@wails dev -v 2 -loglevel Debug 2>&1 | tee wails-dev.log

# Check if dev server processes are running
dev-gui-status:
	@echo "Checking dev server status..."
	@echo ""
	@echo "Running processes:"
	@ps aux | grep -E "(wails|vite|bun.*dev)" | grep -v grep || echo "  No dev server processes found"
	@echo ""
	@echo "Checking if Vite is responding on port 5173:"
	@curl -s -o /dev/null -w "  HTTP Status: %{http_code}\n" http://localhost:5173 2>&1 || echo "  Vite server not responding"
	@echo ""
	@echo "To kill stuck processes, run: make dev-gui-kill"

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
	@echo "    make dev-gui-verbose - Run Wails dev server with verbose output"
	@echo "    make dev-gui-status - Check if dev server is running"
	@echo "    make dev-gui-kill   - Kill stuck dev server processes"
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
