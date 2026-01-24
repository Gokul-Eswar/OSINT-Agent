@echo off
setlocal

:: Set Project Root to the directory of this script
set "PROJECT_ROOT=%~dp0"

:: Set path to the local virtual environment
set "VENV_PATH=%PROJECT_ROOT%.venv"
set "PATH=%VENV_PATH%\Scripts;%PATH%"

:: Add Project Root to PYTHONPATH so 'analyzer' module is found
set "PYTHONPATH=%PROJECT_ROOT%;%PYTHONPATH%"

:: Run the spectre binary with arguments
"%PROJECT_ROOT%spectre.exe" %*

endlocal