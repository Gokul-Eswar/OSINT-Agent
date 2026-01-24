# Dynamic Setup Script

$ScriptPath = $MyInvocation.MyCommand.Path
$ProjectDir = Split-Path $ScriptPath -Parent

Write-Host "Setting up global access for Spectre..." -ForegroundColor Cyan
Write-Host "Project Directory: $ProjectDir" -ForegroundColor Gray

# Get current User PATH
$CurrentPath = [Environment]::GetEnvironmentVariable("Path", "User")

# Normalize paths for comparison (remove trailing slashes, lowercase)
$NormProjectDir = $ProjectDir.TrimEnd('\').ToLower()
$NormCurrentPath = $CurrentPath.ToLower()

if ($NormCurrentPath -like "*$NormProjectDir*") {
    Write-Host "Spectre is already in your PATH." -ForegroundColor Yellow
} else {
    # Add to PATH
    $NewPath = "$CurrentPath;$ProjectDir"
    [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
    
    Write-Host "Success! Added Spectre to your User PATH." -ForegroundColor Green
    Write-Host "Please RESTART your terminal/shell to use the 'spectre' command anywhere." -ForegroundColor Yellow
}