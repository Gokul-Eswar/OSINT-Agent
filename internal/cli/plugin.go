package cli

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Gokul-Eswar/Spectre/internal/collector"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Extension struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Version     string   `json:"version"`
	Author      string   `json:"author"`
	Type        string   `json:"type"`
	Tags        []string `json:"tags"`
}

type Registry struct {
	Extensions []Extension `json:"extensions"`
}

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage external plugins and extensions",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return ensureCollectorBootstrap()
	},
}

var installPluginCmd = &cobra.Command{
	Use:   "install [url|name]",
	Short: "Install a plugin from a URL or registry name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]

		// 1. Resolve target (is it a URL or a name?)
		url := target
		if !strings.HasPrefix(target, "http") {
			fmt.Printf("🔍 Searching registry for plugin '%s'...\n", target)
			reg, err := fetchRegistry()
			if err != nil {
				return err
			}

			found := false
			for _, ext := range reg.Extensions {
				if ext.Name == target {
					return installExtension(ext)
				}
			}

			if !found {
				return fmt.Errorf("plugin '%s' not found in registry", target)
			}
		}

		// Direct URL install
		// Handle GitHub URLs (convert to zip)
		if strings.Contains(url, "github.com") && !strings.HasSuffix(url, ".zip") {
			url = strings.TrimSuffix(url, "/") + "/archive/refs/heads/main.zip"
		}

		fmt.Printf("📥 Installing plugin from %s...\n", url)
		return downloadAndExtract(url)
	},
}

var updatePluginCmd = &cobra.Command{
	Use:   "update [name|all]",
	Short: "Update installed plugins to the latest version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]

		entries, err := os.ReadDir("plugins")
		if err != nil {
			return err
		}

		reg, err := fetchRegistry()
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				if target != "all" && entry.Name() != target {
					continue
				}

				metadata, err := readPluginMetadata(entry.Name())
				if err != nil {
					continue
				}

				for _, ext := range reg.Extensions {
					if ext.Name == metadata.Name {
						if ext.Version != metadata.Version {
							fmt.Printf("🆙 Updating %s from %s to %s...\n", ext.Name, metadata.Version, ext.Version)
							if err := installExtension(ext); err != nil {
								fmt.Printf("❌ Failed to update %s: %v\n", ext.Name, err)
							}
						} else if target != "all" {
							fmt.Printf("✅ %s is already up to date (%s).\n", ext.Name, metadata.Version)
						}
						break
					}
				}
			}
		}
		return nil
	},
}

func installExtension(ext Extension) error {
	url := ext.URL
	// Handle GitHub URLs (convert to zip)
	if strings.Contains(url, "github.com") && !strings.HasSuffix(url, ".zip") {
		url = strings.TrimSuffix(url, "/") + "/archive/refs/heads/main.zip"
	}

	fmt.Printf("📥 Downloading %s (%s)...\n", ext.Name, ext.Version)
	return downloadAndExtract(url)
}

var searchPluginCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search the remote registry for plugins",
	RunE: func(cmd *cobra.Command, args []string) error {
		query := ""
		if len(args) > 0 {
			query = strings.ToLower(args[0])
		}

		reg, err := fetchRegistry()
		if err != nil {
			return err
		}

		fmt.Println("AVAILABLE PLUGINS:")
		fmt.Printf("%-20s %-10s %-15s %-40s\n", "NAME", "VERSION", "TYPE", "DESCRIPTION")
		fmt.Println(strings.Repeat("─", 85))

		for _, ext := range reg.Extensions {
			match := query == "" ||
				strings.Contains(strings.ToLower(ext.Name), query) ||
				strings.Contains(strings.ToLower(ext.Description), query)

			if !match {
				for _, tag := range ext.Tags {
					if strings.Contains(strings.ToLower(tag), query) {
						match = true
						break
					}
				}
			}

			if match {
				fmt.Printf("%-20s %-10s %-15s %-40s\n", ext.Name, ext.Version, ext.Type, ext.Description)
			}
		}
		return nil
	},
}

var infoPluginCmd = &cobra.Command{
	Use:   "info [name]",
	Short: "Show detailed information about a plugin",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		reg, err := fetchRegistry()
		if err != nil {
			return err
		}

		for _, ext := range reg.Extensions {
			if ext.Name == name {
				fmt.Printf("Plugin:      %s\n", ext.Name)
				fmt.Printf("Version:     %s\n", ext.Version)
				fmt.Printf("Author:      %s\n", ext.Author)
				fmt.Printf("Type:        %s\n", ext.Type)
				fmt.Printf("Tags:        %s\n", strings.Join(ext.Tags, ", "))
				fmt.Printf("Description: %s\n", ext.Description)
				fmt.Printf("URL:         %s\n", ext.URL)
				return nil
			}
		}

		return fmt.Errorf("plugin '%s' not found", name)
	},
}

var checkUpdatesPluginCmd = &cobra.Command{
	Use:   "check-updates",
	Short: "Check for updates to installed plugins",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := os.ReadDir("plugins")
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("No external plugins installed.")
				return nil
			}
			return err
		}

		reg, err := fetchRegistry()
		if err != nil {
			return err
		}

		fmt.Println("CHECKING FOR UPDATES:")
		fmt.Printf("%-20s %-15s %-15s %-10s\n", "NAME", "INSTALLED", "LATEST", "STATUS")
		fmt.Println(strings.Repeat("─", 65))

		for _, entry := range entries {
			if entry.IsDir() {
				metadata, err := readPluginMetadata(entry.Name())
				if err != nil {
					continue
				}

				latestVersion := "unknown"
				status := "up to date"

				for _, ext := range reg.Extensions {
					if ext.Name == metadata.Name {
						latestVersion = ext.Version
						if latestVersion != metadata.Version {
							status = "UPDATE AVAILABLE"
						}
						break
					}
				}

				fmt.Printf("%-20s %-15s %-15s %-10s\n", metadata.Name, metadata.Version, latestVersion, status)
			}
		}
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
		fmt.Printf("%-20s %-10s %-40s\n", "NAME", "VERSION", "DESCRIPTION")
		fmt.Println(strings.Repeat("─", 75))

		found := false
		for _, entry := range entries {
			if entry.IsDir() {
				metadata, err := readPluginMetadata(entry.Name())
				if err != nil {
					fmt.Printf("- %s (error reading metadata)\n", entry.Name())
					continue
				}
				fmt.Printf("%-20s %-10s %-40s\n", metadata.Name, metadata.Version, metadata.Description)
				found = true
			}
		}
		if !found {
			fmt.Println("No external plugins found in the plugins/ directory.")
		}
		return nil
	},
}

func readPluginMetadata(pluginName string) (*collector.PluginMetadata, error) {
	path := filepath.Join("plugins", pluginName, "plugin.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var meta collector.PluginMetadata
	if err := yaml.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	return &meta, nil
}

func fetchRegistry() (*Registry, error) {
	url := viper.GetString("extension_registry_url")
	if url == "" {
		url = "file://registry_sample.json" // Fallback for local dev
	}

	var data []byte
	var err error

	if strings.HasPrefix(url, "file://") {
		path := strings.TrimPrefix(url, "file://")
		data, err = os.ReadFile(path)
	} else {
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch registry: %w", err)
		}
		defer resp.Body.Close()
		data, err = io.ReadAll(resp.Body)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read registry data: %w", err)
	}

	var reg Registry
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, fmt.Errorf("failed to parse registry JSON: %w", err)
	}

	return &reg, nil
}

func downloadAndExtract(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: status %d", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "spectre-plugin-*.zip")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err = io.Copy(tmpFile, resp.Body); err != nil {
		return err
	}

	pluginsDir := "plugins"
	os.MkdirAll(pluginsDir, 0755)

	r, err := zip.OpenReader(tmpFile.Name())
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(pluginsDir, f.Name)
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

	fmt.Println("Plugin installed successfully.")
	return nil
}

func init() {
	pluginCmd.AddCommand(installPluginCmd)
	pluginCmd.AddCommand(searchPluginCmd)
	pluginCmd.AddCommand(infoPluginCmd)
	pluginCmd.AddCommand(checkUpdatesPluginCmd)
	pluginCmd.AddCommand(updatePluginCmd)
	pluginCmd.AddCommand(listPluginsCmd)
	rootCmd.AddCommand(pluginCmd)
}
