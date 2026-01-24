@echo off
setlocal

:: Set path to the local virtual environment
set "VENV_PATH=%~dp0.venv"
set "PATH=%VENV_PATH%\Scripts;%PATH%"

:: Run the spectre core binary with arguments
"%~dp0spectre-core.exe" %*

endlocal