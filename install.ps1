# GitR Installation Script for Windows
# This script installs GitR system-wide and makes it available as 'gitr' command

param(
    [switch]$Force
)

# Colors for output
$Red = "Red"
$Green = "Green"
$Yellow = "Yellow"
$Blue = "Blue"

# Function to print colored output
function Write-Status {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor $Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor $Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor $Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor $Red
}

# Check if Go is installed
function Test-Go {
    try {
        $goVersion = go version 2>$null
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Go is installed: $goVersion"
            return $true
        }
    }
    catch {
        Write-Error "Go is not installed. Please install Go 1.21 or later."
        Write-Status "Visit: https://golang.org/doc/install"
        return $false
    }
}

# Build GitR
function Build-GitR {
    Write-Status "Building GitR..."
    
    if (-not (Test-Path "main.go")) {
        Write-Error "main.go not found. Please run this script from the GitR source directory."
        exit 1
    }
    
    # Build the binary
    go build -o gitr.exe .
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "GitR built successfully"
    } else {
        Write-Error "Failed to build GitR"
        exit 1
    }
}

# Install GitR system-wide
function Install-GitR {
    Write-Status "Installing GitR system-wide..."
    
    # Determine installation directory
    $installDir = "$env:ProgramFiles\GitR"
    $binDir = "$installDir\bin"
    
    # Create installation directory
    if (-not (Test-Path $installDir)) {
        New-Item -ItemType Directory -Path $installDir -Force | Out-Null
    }
    
    if (-not (Test-Path $binDir)) {
        New-Item -ItemType Directory -Path $binDir -Force | Out-Null
    }
    
    # Copy binary
    Copy-Item "gitr.exe" "$binDir\gitr.exe" -Force
    
    # Add to PATH if not already present
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "Machine")
    if ($currentPath -notlike "*$binDir*") {
        Write-Status "Adding GitR to system PATH..."
        $newPath = "$currentPath;$binDir"
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "Machine")
        Write-Success "GitR added to system PATH"
    } else {
        Write-Success "GitR already in system PATH"
    }
    
    # Verify installation
    try {
        $gitrPath = Get-Command gitr -ErrorAction Stop
        Write-Success "GitR installed successfully at $($gitrPath.Source)"
    }
    catch {
        Write-Warning "GitR installed but may not be available in current session."
        Write-Status "Please restart your terminal or run: refreshenv"
    }
}

# Clean up build artifacts
function Remove-BuildArtifacts {
    Write-Status "Cleaning up build artifacts..."
    if (Test-Path "gitr.exe") {
        Remove-Item "gitr.exe" -Force
    }
    Write-Success "Cleanup completed"
}

# Main installation process
function Main {
    Write-Host "=========================================="
    Write-Host "        GitR Installation Script"
    Write-Host "=========================================="
    Write-Host ""
    
    Write-Status "Starting GitR installation..."
    
    # Check if running as administrator
    $isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")
    
    if (-not $isAdmin -and -not $Force) {
        Write-Warning "This script requires administrator privileges to install system-wide."
        Write-Status "Please run PowerShell as Administrator, or use -Force flag for user installation."
        exit 1
    }
    
    # Check prerequisites
    if (-not (Test-Go)) {
        exit 1
    }
    
    # Build GitR
    Build-GitR
    
    # Install GitR
    Install-GitR
    
    # Clean up
    Remove-BuildArtifacts
    
    Write-Host ""
    Write-Host "=========================================="
    Write-Success "GitR installation completed!"
    Write-Host "=========================================="
    Write-Host ""
    Write-Status "You can now use 'gitr' command from anywhere."
    Write-Status "Run 'gitr --help' to see available options."
    Write-Status "Run 'gitr' in a git repository to get started."
    Write-Host ""
}

# Run main function
Main
