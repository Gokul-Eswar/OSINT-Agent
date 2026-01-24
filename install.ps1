Write-Host "SPECTRE System Installer" -ForegroundColor Cyan
Write-Host "========================" -ForegroundColor Cyan

# 1. Check Go
if (Get-Command go -ErrorAction SilentlyContinue) {
    Write-Host "[+] Go is installed" -ForegroundColor Green
} else {
    Write-Host "[-] Go is NOT installed. Please install Go 1.21+." -ForegroundColor Red
    exit 1
}

# 2. Check Python
if (Get-Command python -ErrorAction SilentlyContinue) {
    Write-Host "[+] Python is installed" -ForegroundColor Green
} else {
    Write-Host "[-] Python is NOT installed. Please install Python 3.10+." -ForegroundColor Red
    exit 1
}

# 3. Build Go Binary
Write-Host "[*] Building SPECTRE binary..." -ForegroundColor Yellow
go build -o spectre-core.exe cmd/spectre/main.go
if ($LastExitCode -eq 0) {
    Write-Host "[+] Build successful: spectre-core.exe" -ForegroundColor Green
} else {
    Write-Host "[-] Build failed." -ForegroundColor Red
    exit 1
}

# 4. Setup Python Venv
$VenvDir = ".venv"
if (-Not (Test-Path $VenvDir)) {
    Write-Host "[*] Creating Python virtual environment..." -ForegroundColor Yellow
    python -m venv $VenvDir
} else {
    Write-Host "[*] Python virtual environment exists." -ForegroundColor Gray
}

# 5. Install Requirements
Write-Host "[*] Installing Python dependencies..." -ForegroundColor Yellow
& ".\$VenvDir\Scripts\python.exe" -m pip install -r analyzer/requirements.txt
if ($LastExitCode -eq 0) {
    Write-Host "[+] Dependencies installed." -ForegroundColor Green
} else {
    Write-Host "[-] Dependency installation failed." -ForegroundColor Red
    exit 1
}

Write-Host "`n[+] Installation Complete!" -ForegroundColor Green
Write-Host "Run '.\spectre.bat' to start the system." -ForegroundColor Cyan
