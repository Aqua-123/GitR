#!/bin/bash

# GitR Installation Script
# This script installs GitR system-wide and makes it available as 'gitr' command

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.21 or later."
        print_status "Visit: https://golang.org/doc/install"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_success "Go is installed: version $GO_VERSION"
}

# Build GitR
build_gitr() {
    print_status "Building GitR..."
    
    if [ ! -f "main.go" ]; then
        print_error "main.go not found. Please run this script from the GitR source directory."
        exit 1
    fi
    
    # Build the binary
    go build -o gitr .
    
    if [ $? -eq 0 ]; then
        print_success "GitR built successfully"
    else
        print_error "Failed to build GitR"
        exit 1
    fi
}

# Install GitR system-wide
install_gitr() {
    print_status "Installing GitR system-wide..."
    
    # Determine installation directory
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        INSTALL_DIR="/usr/local/bin"
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        INSTALL_DIR="/usr/local/bin"
    else
        print_error "Unsupported operating system: $OSTYPE"
        exit 1
    fi
    
    # Check if we have write permissions
    if [ ! -w "$INSTALL_DIR" ]; then
        print_warning "No write permission to $INSTALL_DIR. Using sudo..."
        sudo cp gitr "$INSTALL_DIR/"
        sudo chmod +x "$INSTALL_DIR/gitr"
    else
        cp gitr "$INSTALL_DIR/"
        chmod +x "$INSTALL_DIR/gitr"
    fi
    
    # Verify installation
    if command -v gitr &> /dev/null; then
        print_success "GitR installed successfully at $(which gitr)"
    else
        print_error "Installation failed. GitR command not found in PATH."
        exit 1
    fi
}

# Clean up build artifacts
cleanup() {
    print_status "Cleaning up build artifacts..."
    rm -f gitr
    print_success "Cleanup completed"
}

# Main installation process
main() {
    echo "=========================================="
    echo "        GitR Installation Script"
    echo "=========================================="
    echo ""
    
    print_status "Starting GitR installation..."
    
    # Check prerequisites
    check_go
    
    # Build GitR
    build_gitr
    
    # Install GitR
    install_gitr
    
    # Clean up
    cleanup
    
    echo ""
    echo "=========================================="
    print_success "GitR installation completed!"
    echo "=========================================="
    echo ""
    print_status "You can now use 'gitr' command from anywhere."
    print_status "Run 'gitr --help' to see available options."
    print_status "Run 'gitr' in a git repository to get started."
    echo ""
}

# Run main function
main "$@"
