package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spectre/spectre/internal/extensions"
)

var extCmd = &cobra.Command{
	Use:   "extension",
	Short: "Manage extensions (search, install, list)",
	Long:  `Search for and install extensions from the official Spectre registry.`,
	Aliases: []string{"ext"},
}

var extListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed extensions",
	Run: func(cmd *cobra.Command, args []string) {
		mgr := extensions.NewManager("plugins")
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
		for _, name := range installed {
			fmt.Printf("- %s\n", name)
		}
	},
}

var extSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search the extension registry",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mgr := extensions.NewManager("plugins")
		
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

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tDESCRIPTION\tAUTHOR\tVERSION")
		for _, ext := range results {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", ext.Name, ext.Description, ext.Author, ext.Version)
		}
		w.Flush()
	},
}

var extInstallCmd = &cobra.Command{
	Use:   "install [name]",
	Short: "Install an extension",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mgr := extensions.NewManager("plugins")
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
		mgr := extensions.NewManager("plugins")
		err := mgr.Remove(args[0])
		if err != nil {
			fmt.Printf("Error removing extension: %v\n", err)
			return
		}
		fmt.Printf("Successfully removed '%s'.\n", args[0])
	},
}

func init() {
	extCmd.AddCommand(extListCmd)
	extCmd.AddCommand(extSearchCmd)
	extCmd.AddCommand(extInstallCmd)
	extCmd.AddCommand(extRemoveCmd)
	
	// Register with root
	rootCmd.AddCommand(extCmd)
}