# scripts/build_analyzer.ps1
$ErrorActionPreference = "Stop"

Write-Host "==================================================" -ForegroundColor Cyan
Write-Host "     Freezing Python Analyzer with PyInstaller    " -ForegroundColor Cyan
Write-Host "==================================================" -ForegroundColor Cyan

# Verify PyInstaller installation
if (-not (Get-Command "pyinstaller" -ErrorAction SilentlyContinue)) {
    Write-Host "[*] PyInstaller not found. Installing..." -ForegroundColor Yellow
    pip install pyinstaller -q
}

# Run PyInstaller using spec file
Write-Host "[*] Compiling standalone analyzer bundle into dist/spectre-analyzer..." -ForegroundColor Cyan
pyinstaller --noconfirm scripts/analyzer.spec --distpath dist

if ($LASTEXITCODE -eq 0) {
    Write-Host "==================================================" -ForegroundColor Green
    Write-Host " ✅ Analyzer frozen successfully at: dist/spectre-analyzer" -ForegroundColor Green
    Write-Host "==================================================" -ForegroundColor Green
} else {
    Write-Error "PyInstaller build failed."
    exit 1
}
