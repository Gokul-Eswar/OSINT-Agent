$ProjectDir = "C:\Users\GOKUL\Desktop\CLI_AGENNT"

# Get current User PATH
$CurrentPath = [Environment]::GetEnvironmentVariable("Path", "User")

# Check if already in PATH
if ($CurrentPath -like "*$ProjectDir*") {
    Write-Host "Spectre is already in your PATH." -ForegroundColor Yellow
} else {
    # Add to PATH
    $NewPath = "$CurrentPath;$ProjectDir"
    [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
    Write-Host "Success! Added Spectre to your User PATH." -ForegroundColor Green
    Write-Host "You must RESTART your terminal for this to take effect." -ForegroundColor Cyan
}
