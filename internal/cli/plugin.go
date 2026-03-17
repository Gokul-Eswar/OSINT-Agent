package cli

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage external plugins",
}

var installPluginCmd = &cobra.Command{
	Use:   "install [url]",
	Short: "Install a plugin from a URL (ZIP format)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]
		fmt.Printf("Installing plugin from %s...\n", url)

		// 1. Download
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to download plugin: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to download plugin: status %d", resp.StatusCode)
		}

		tmpFile, err := os.CreateTemp("", "spectre-plugin-*.zip")
		if err != nil {
			return fmt.Errorf("failed to create temp file: %w", err)
		}
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		_, err = io.Copy(tmpFile, resp.Body)
		if err != nil {
			return fmt.Errorf("failed to save temp file: %w", err)
		}

		// 2. Extract
		pluginsDir := "plugins"
		if err := os.MkdirAll(pluginsDir, 0755); err != nil {
			return fmt.Errorf("failed to create plugins directory: %w", err)
		}

		r, err := zip.OpenReader(tmpFile.Name())
		if err != nil {
			return fmt.Errorf("failed to open zip: %w", err)
		}
		defer r.Close()

		for _, f := range r.File {
			fpath := filepath.Join(pluginsDir, f.Name)

			// Check for ZipSlip
			if !strings.HasPrefix(fpath, filepath.Clean(pluginsDir)+string(os.PathSeparator)) && fpath != pluginsDir {
				return fmt.Errorf("illegal file path: %s", fpath)
			}

			if f.FileInfo().IsDir() {
				os.MkdirAll(fpath, os.ModePerm)
				continue
			}

			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			rc, err := f.Open()
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, rc)

			outFile.Close()
			rc.Close()

			if err != nil {
				return err
			}
		}

		fmt.Println("Plugin installed successfully. It will be loaded on the next run.")
		return nil
	},
}

var listPluginsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all installed external plugins",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := os.ReadDir("plugins")
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("No external plugins installed.")
				return nil
			}
			return err
		}

		fmt.Println("INSTALLED EXTERNAL PLUGINS:")
		fmt.Println("────────────────────────────")
		found := false
		for _, entry := range entries {
			if entry.IsDir() {
				fmt.Printf("- %s\n", entry.Name())
				found = true
			}
		}
		if !found {
			fmt.Println("No external plugins found in the plugins/ directory.")
		}
		return nil
	},
}

func init() {
	pluginCmd.AddCommand(installPluginCmd)
	pluginCmd.AddCommand(listPluginsCmd)
	rootCmd.AddCommand(pluginCmd)
}
