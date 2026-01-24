package cli

import (
	"fmt"
	"os"

	"github.com/spectre/spectre/internal/config"
	"github.com/spectre/spectre/internal/logger"
	"github.com/spf13/cobra"
)

var cfgFile string
var activeAllowed bool

var rootCmd = &cobra.Command{
	Use:   "spectre",
	Short: "SPECTRE - Local-First OSINT Intelligence Platform",
	Long: `SPECTRE turns raw internet noise into structured, verifiable intelligence — fast, repeatable, and local.
Not scraping. Not search. Intelligence synthesis with auditability.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("SPECTRE CLI initialized")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.spectre.yaml)")
	rootCmd.PersistentFlags().BoolVar(&activeAllowed, "active", false, "Allow active reconnaissance (port scans, etc.)")
}

func initConfig() {
	config.InitConfig(cfgFile)
	logger.InitLogger()
}
