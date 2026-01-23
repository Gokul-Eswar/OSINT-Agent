package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage SPECTRE configuration",
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		viper.Set(key, value)
		
		err := viper.WriteConfig()
		if err != nil {
			// If config file not found, try to create it
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				err = viper.SafeWriteConfig()
			}
			if err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
		}

		fmt.Printf("Successfully set %s = %s\n", key, value)
		fmt.Printf("Config saved to: %s\n", viper.ConfigFileUsed())
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		settings := viper.AllSettings()
		fmt.Println("Current Configuration:")
		for k, v := range settings {
			fmt.Printf("%s: %v\n", k, v)
		}
		fmt.Printf("\nConfig file: %s\n", viper.ConfigFileUsed())
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)
	rootCmd.AddCommand(configCmd)
}
