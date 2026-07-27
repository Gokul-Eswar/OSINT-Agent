package main

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed payload.zip
var payloadZip []byte

func main() {
	fmt.Println("==================================================")
	fmt.Println("    🕵️  SPECTRE - Single Executable Installer")
	fmt.Println("==================================================")

	// Determine default installation directory (%LocalAppData%\SPECTRE on Windows)
	defaultDir := os.Getenv("LOCALAPPDATA")
	if defaultDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			defaultDir = "."
		} else {
			defaultDir = filepath.Join(home, "AppData", "Local")
		}
	}
	installDir := filepath.Join(defaultDir, "SPECTRE")

	// Allow command-line flag override
	if len(os.Args) > 1 {
		for i, arg := range os.Args {
			if (arg == "-dir" || arg == "--dir") && i+1 < len(os.Args) {
				installDir = os.Args[i+1]
			}
		}
	}

	fmt.Printf("[*] Installing SPECTRE to: %s\n", installDir)

	// 1. Create installation directory
	if err := os.MkdirAll(installDir, 0755); err != nil {
		fmt.Printf("[-] Failed to create directory %s: %v\n", installDir, err)
		os.Exit(1)
	}

	// 2. Extract embedded payload zip
	fmt.Println("[*] Unpacking application files...")
	if err := extractZip(payloadZip, installDir); err != nil {
		fmt.Printf("[-] Extraction failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("[+] Application files unpacked successfully.")

	// 3. Setup Python virtual environment if python is available
	pythonExe := findPython()
	if pythonExe != "" {
		fmt.Printf("[*] Found Python: %s\n", pythonExe)
		venvDir := filepath.Join(installDir, ".venv")
		venvPython := filepath.Join(venvDir, "Scripts", "python.exe")

		if _, err := os.Stat(venvPython); os.IsNotExist(err) {
			fmt.Println("[*] Creating Python virtual environment...")
			cmd := exec.Command(pythonExe, "-m", "venv", venvDir)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				// Retry with --without-pip if ensurepip failed
				cmdFallback := exec.Command(pythonExe, "-m", "venv", "--without-pip", venvDir)
				if errFb := cmdFallback.Run(); errFb != nil {
					fmt.Printf("[!] Warning: Failed to create virtualenv: %v\n", errFb)
				}
			}
		}

		reqFile := filepath.Join(installDir, "analyzer", "requirements.txt")
		if _, err := os.Stat(venvPython); err == nil {
			if _, err := os.Stat(reqFile); err == nil {
				fmt.Println("[*] Installing Python dependencies...")
				cmd := exec.Command(venvPython, "-m", "pip", "install", "-r", reqFile)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					fmt.Printf("[!] Warning: Dependency installation failed: %v\n", err)
				} else {
					fmt.Println("[+] Python dependencies installed.")
				}
			}
		}
	} else {
		fmt.Println("[!] Python not found in system PATH. Install Python 3.10+ to enable AI features.")
	}

	// 4. Add installDir to User PATH environment variable
	fmt.Println("[*] Configuring PATH environment variable...")
	if err := addToUserPath(installDir); err != nil {
		fmt.Printf("[!] Warning: Could not update PATH automatically: %v\n", err)
		fmt.Printf("    You can manually add '%s' to your User PATH environment variable.\n", installDir)
	} else {
		fmt.Println("[+] PATH updated successfully.")
	}

	fmt.Println("\n==================================================")
	fmt.Println("    ✅ SPECTRE Installation Complete!")
	fmt.Println("==================================================")
	fmt.Printf(" Installed location: %s\n", installDir)
	fmt.Println(" Run 'spectre' or 'spectre.bat' from any new terminal window to start.")
	fmt.Println("==================================================")
}

func extractZip(zipData []byte, targetDir string) error {
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return err
	}

	for _, file := range zipReader.File {
		path := filepath.Join(targetDir, file.Name)

		// Guard against zip slip security issues
		if !strings.HasPrefix(filepath.Clean(path), filepath.Clean(targetDir)) {
			return fmt.Errorf("illegal file path: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func findPython() string {
	for _, py := range []string{"python", "python3", "py"} {
		if path, err := exec.LookPath(py); err == nil {
			out, err := exec.Command(path, "--version").Output()
			if err == nil && len(out) > 0 {
				return path
			}
		}
	}
	return ""
}

func addToUserPath(dir string) error {
	// Execute PowerShell command to add directory to User PATH environment variable
	psCmd := fmt.Sprintf(`
$dir = '%s'
$oldPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if ($oldPath -notlike "*$dir*") {
    $newPath = "$oldPath;$dir"
    [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
}
`, dir)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psCmd)
	return cmd.Run()
}
