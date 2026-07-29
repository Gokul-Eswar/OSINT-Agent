package cli

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Gokul-Eswar/Spectre/internal/extensions"
	"github.com/Gokul-Eswar/Spectre/internal/tui/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var extCmd = &cobra.Command{
	Use:     "extension",
	Short:   "Manage extensions (search, install, list)",
	Long:    `Search for and install extensions from the official Spectre registry.`,
	Aliases: []string{"ext"},
	Run: func(cmd *cobra.Command, args []string) {
		// Default to UI
		startTUI()
	},
}

func getManager() *extensions.Manager {
	url := viper.GetString("extension_registry_url")
	if url == "" {
		// Fallback if not configured
		url = "https://raw.githubusercontent.com/spectre-org/registry/main/registry.json"
	}
	return extensions.NewManager("plugins", url)
}

func startTUI() {
	mgr := getManager()
	model := store.NewModel(mgr)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error starting store UI: %v\n", err)
		os.Exit(1)
	}
}

var extUiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Open the interactive extension store",
	Run: func(cmd *cobra.Command, args []string) {
		startTUI()
	},
}

var tagFilter string

var extListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed extensions",
	Run: func(cmd *cobra.Command, args []string) {
		mgr := getManager()
		installed, err := mgr.ListInstalled()
		if err != nil {
			fmt.Printf("Error listing extensions: %v\n", err)
			return
		}

		if len(installed) == 0 {
			fmt.Println("No extensions installed.")
			return
		}

		fmt.Println("Installed Extensions:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tDESCRIPTION\tTAGS")
		for _, ext := range installed {
			tags := strings.Join(ext.Tags, ", ")
			if tags == "" {
				tags = "-"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", ext.Name, ext.Description, tags)
		}
		w.Flush()
	},
}

var extSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search the extension registry",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mgr := getManager()

		query := ""
		if len(args) > 0 {
			query = args[0]
		}

		results, err := mgr.Search(query)
		if err != nil {
			fmt.Printf("Error searching: %v\n", err)
			return
		}

		if len(results) == 0 {
			fmt.Println("No extensions found matching your query.")
			return
		}

		// Filter by tag if provided
		if tagFilter != "" {
			var filtered []extensions.Extension
			for _, ext := range results {
				match := false
				for _, t := range ext.Tags {
					if strings.EqualFold(t, tagFilter) {
						match = true
						break
					}
				}
				if match {
					filtered = append(filtered, ext)
				}
			}
			results = filtered
		}

		if len(results) == 0 {
			fmt.Printf("No extensions found with tag '%s'.\n", tagFilter)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tDESCRIPTION\tTAGS\tVERSION")
		for _, ext := range results {
			tags := strings.Join(ext.Tags, ", ")
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", ext.Name, ext.Description, tags, ext.Version)
		}
		w.Flush()
	},
}

var extInstallCmd = &cobra.Command{
	Use:   "install [name]",
	Short: "Install an extension",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mgr := getManager()
		err := mgr.Install(args[0])
		if err != nil {
			fmt.Printf("Error installing extension: %v\n", err)
			return
		}
		fmt.Printf("Successfully installed '%s'. It is now available for use.\n", args[0])
	},
}

var extRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove an extension",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mgr := getManager()
		err := mgr.Remove(args[0])
		if err != nil {
			fmt.Printf("Error removing extension: %v\n", err)
			return
		}
		fmt.Printf("Successfully removed '%s'.\n", args[0])
	},
}

var extUpdateCmd = &cobra.Command{
	Use:   "update [name|all]",
	Short: "Update one or all installed extensions",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mgr := getManager()
		if args[0] == "all" {
			err := mgr.UpdateAll()
			if err != nil {
				fmt.Printf("Error updating extensions: %v\n", err)
			}
		} else {
			err := mgr.Update(args[0])
			if err != nil {
				fmt.Printf("Error updating extension: %v\n", err)
			}
		}
	},
}

var extInfoCmd = &cobra.Command{
	Use:   "info [name]",
	Short: "Show detailed information about an extension",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mgr := getManager()
		ext, installed, err := mgr.GetInfo(args[0])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Name:        %s\n", ext.Name)
		fmt.Printf("Description: %s\n", ext.Description)
		fmt.Printf("Author:      %s\n", ext.Author)
		fmt.Printf("Version:     %s\n", ext.Version)
		fmt.Printf("Type:        %s\n", ext.Type)
		fmt.Printf("Tags:        %s\n", strings.Join(ext.Tags, ", "))
		fmt.Printf("URL:         %s\n", ext.URL)

		status := "Not installed"
		if installed {
			status = "Installed"
		}
		fmt.Printf("Status:      %s\n", status)
	},
}

func init() {
	extSearchCmd.Flags().StringVarP(&tagFilter, "tag", "t", "", "Filter extensions by tag")

	extCmd.AddCommand(extListCmd)

	extCmd.AddCommand(extSearchCmd)
	extCmd.AddCommand(extInstallCmd)
	extCmd.AddCommand(extRemoveCmd)
	extCmd.AddCommand(extUpdateCmd)
	extCmd.AddCommand(extInfoCmd)
	extCmd.AddCommand(extUiCmd)

	// Register with root
	rootCmd.AddCommand(extCmd)
}
