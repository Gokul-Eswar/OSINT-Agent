# scripts/build_release.ps1

$Version = "v1.0.0"
$DistDir = "dist"

# Clean dist directory
if (Test-Path $DistDir) {
    Remove-Item -Recurse -Force $DistDir
}
New-Item -ItemType Directory -Force -Path $DistDir | Out-Null

Write-Host "Building SPECTRE $Version..." -ForegroundColor Cyan

# Windows amd64
Write-Host "Compiling for Windows (amd64)..."
$Env:GOOS = "windows"
$Env:GOARCH = "amd64"
go build -o "$DistDir/spectre_${Version}_windows_amd64.exe" ./cmd/spectre
if ($LASTEXITCODE -ne 0) { Write-Error "Windows build failed"; exit 1 }

# Linux amd64
Write-Host "Compiling for Linux (amd64)..."
$Env:GOOS = "linux"
$Env:GOARCH = "amd64"
go build -o "$DistDir/spectre_${Version}_linux_amd64" ./cmd/spectre
if ($LASTEXITCODE -ne 0) { Write-Error "Linux build failed"; exit 1 }

# Darwin arm64 (macOS Apple Silicon)
Write-Host "Compiling for macOS (arm64)..."
$Env:GOOS = "darwin"
$Env:GOARCH = "arm64"
go build -o "$DistDir/spectre_${Version}_darwin_arm64" ./cmd/spectre
if ($LASTEXITCODE -ne 0) { Write-Error "macOS build failed"; exit 1 }

Write-Host "Build complete! Artifacts in $DistDir/" -ForegroundColor Green
