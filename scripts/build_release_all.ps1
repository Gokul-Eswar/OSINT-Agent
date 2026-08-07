# scripts/build_release_all.ps1
# Master script to compile Go binary, freeze Python analyzer (--onedir), and build single-click installer

$ErrorActionPreference = "Stop"

Write-Host "==================================================" -ForegroundColor Cyan
Write-Host " SPECTRE Zero-Dependency Master Build Pipeline " -ForegroundColor Cyan
Write-Host "==================================================" -ForegroundColor Cyan

# 1. Build main Go executable statically
Write-Host "[*] Step 1: Compiling main SPECTRE binary (spectre.exe)..." -ForegroundColor Yellow
go build -ldflags="-s -w" -o spectre.exe ./cmd/spectre
if ($LASTEXITCODE -ne 0) { Write-Error "Failed to build spectre.exe"; exit 1 }
Write-Host "[+] spectre.exe compiled." -ForegroundColor Green

# 2. Freeze Python Analyzer into dist/spectre-analyzer directory bundle
Write-Host "[*] Step 2: Freezing Python Analyzer..." -ForegroundColor Yellow
powershell -File ./scripts/build_analyzer.ps1
if ($LASTEXITCODE -ne 0) { Write-Error "Failed to freeze analyzer"; exit 1 }

# 3. Create Payload Zip for standalone installer
Write-Host "[*] Step 3: Packaging standalone release bundle..." -ForegroundColor Yellow
$PayloadZip = "cmd/installer/payload.zip"
if (Test-Path $PayloadZip) { Remove-Item -Force $PayloadZip }

# Items included in single-click payload
$ItemsToCompress = @(
    "spectre.exe",
    "spectre.bat",
    "dist/spectre-analyzer",
    "configs",
    "web",
    "LICENSE",
    "README.md"
)

Compress-Archive -Path $ItemsToCompress -DestinationPath $PayloadZip -Force
Write-Host "[+] Payload zipped at: $PayloadZip" -ForegroundColor Green

# 4. Build single-click installer executable
Write-Host "[*] Step 4: Compiling zero-dependency installer (spectre-installer.exe)..." -ForegroundColor Yellow
go build -ldflags="-s -w" -o spectre-installer.exe ./cmd/installer
if ($LASTEXITCODE -ne 0) { Write-Error "Failed to build spectre-installer.exe"; exit 1 }

# Clean up payload zip
Remove-Item -Force $PayloadZip

$InstallerSize = (Get-Item spectre-installer.exe).Length / 1MB
Write-Host "`n==================================================" -ForegroundColor Green
Write-Host " 🎉 BUILD SUCCESSFUL!" -ForegroundColor Green
Write-Host " Output: spectre-installer.exe ({0:N2} MB)" -f $InstallerSize -ForegroundColor Cyan
Write-Host " Zero-dependency, 1-click execution ready!" -ForegroundColor Green
Write-Host "==================================================" -ForegroundColor Green
