# scripts/build_installer.ps1
# Script to build a single standalone installer executable (spectre-installer.exe)

$ErrorActionPreference = "Stop"

Write-Host "==================================================" -ForegroundColor Cyan
Write-Host "   Building SPECTRE Single Executable Installer   " -ForegroundColor Cyan
Write-Host "==================================================" -ForegroundColor Cyan

# 1. Build main spectre binary
Write-Host "[*] Compiling spectre.exe..." -ForegroundColor Yellow
go build -o spectre.exe ./cmd/spectre
if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to compile spectre.exe"
    exit 1
}

# 2. Package payload zip for installer embedding
$PayloadZip = "cmd/installer/payload.zip"
if (Test-Path $PayloadZip) {
    Remove-Item -Force $PayloadZip
}

Write-Host "[*] Packaging payload archive..." -ForegroundColor Yellow
$ItemsToCompress = @(
    "spectre.exe",
    "spectre.bat",
    "analyzer",
    "configs",
    "web",
    "LICENSE",
    "README.md"
)

Compress-Archive -Path $ItemsToCompress -DestinationPath $PayloadZip -Force
Write-Host "[+] Payload compressed to $PayloadZip" -ForegroundColor Green

# 3. Build the installer binary
Write-Host "[*] Compiling spectre-installer.exe..." -ForegroundColor Yellow
go build -o spectre-installer.exe ./cmd/installer
if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to compile installer executable"
    exit 1
}

# 4. Clean up temporary payload zip
Remove-Item -Force $PayloadZip

Write-Host "`n==================================================" -ForegroundColor Green
Write-Host "   ✅ spectre-installer.exe created successfully!   " -ForegroundColor Green
Write-Host "==================================================" -ForegroundColor Green
$FileSize = (Get-Item spectre-installer.exe).Length / 1MB
Write-Host (" Output: spectre-installer.exe ({0:N2} MB)" -f $FileSize) -ForegroundColor Cyan
Write-Host " Anyone can run '.\spectre-installer.exe' to install SPECTRE with one click!" -ForegroundColor Yellow
