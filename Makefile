# GitR Makefile
# Provides easy installation and build targets

.PHONY: build install clean help test

# Default target
all: build

# Build GitR binary
build:
	@echo "Building GitR..."
	go build -o gitr .
	@echo "Build completed: gitr"

# Install GitR system-wide (Unix/Linux/macOS)
install: build
	@echo "Installing GitR system-wide..."
	@if [ -w /usr/local/bin ]; then \
		cp gitr /usr/local/bin/ && chmod +x /usr/local/bin/gitr; \
	else \
		sudo cp gitr /usr/local/bin/ && sudo chmod +x /usr/local/bin/gitr; \
	fi
	@echo "GitR installed successfully"

# Install for current user only (Unix/Linux/macOS)
install-user: build
	@echo "Installing GitR for current user..."
	@mkdir -p ~/bin
	@cp gitr ~/bin/
	@chmod +x ~/bin/gitr
	@echo "GitR installed to ~/bin/gitr"
	@echo "Make sure ~/bin is in your PATH"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f gitr gitr.exe
	@echo "Cleanup completed"

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Show help
help:
	@echo "GitR Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  build        - Build GitR binary"
	@echo "  install      - Install GitR system-wide (requires sudo)"
	@echo "  install-user - Install GitR for current user only"
	@echo "  clean        - Remove build artifacts"
	@echo "  test         - Run tests"
	@echo "  help         - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build        # Build the binary"
	@echo "  make install      # Install system-wide"
	@echo "  make install-user # Install for current user"
	@echo "  make clean        # Clean up"

