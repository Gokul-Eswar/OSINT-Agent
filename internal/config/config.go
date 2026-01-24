package config

import (
	"fmt"
	"os"

	"github.com/spectre/spectre/internal/ethics"
	"github.com/spf13/viper"
)

// InitConfig reads in config file and ENV variables if set.
func InitConfig(cfgFile string) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".spectre" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.AddConfigPath("configs") // Added to find default.yaml if present
		viper.SetConfigType("yaml")
		viper.SetConfigName("default") // Default to default.yaml
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	ApplyEthicsConfig()
}

// ApplyEthicsConfig loads settings from viper into the ethics package.
func ApplyEthicsConfig() {
	// Apply Blacklist
	bl := viper.GetStringSlice("ethics.blacklist")
	if len(bl) > 0 {
		ethics.SetBlacklist(bl)
	}

	// Apply Whitelist
	wl := viper.GetStringSlice("ethics.whitelist")
	if len(wl) > 0 {
		ethics.SetWhitelist(wl)
	}

	// Apply Rate Limits
	// We check for collectors.<name>.rate_limit
	collectors := []string{"dns", "whois", "github"}
	for _, name := range collectors {
		key := fmt.Sprintf("collectors.%s.rate_limit", name)
		if viper.IsSet(key) {
			limit := viper.GetFloat64(key)
			ethics.SetLimit(name, limit)
		}
	}
}

// GetAPIKey retrieves an API key from configuration or environment.
func GetAPIKey(name string) string {
	// Checks keys.<name> in config or SPECTRE_KEYS_<NAME> in env
	return viper.GetString("keys." + name)
}
