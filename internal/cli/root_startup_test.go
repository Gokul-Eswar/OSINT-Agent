package cli

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestShouldSkipStartupInit_CommandTiering(t *testing.T) {
	version := &cobra.Command{Use: "version"}
	if !shouldSkipStartupInit(version) {
		t.Fatal("expected version command to skip startup init")
	}

	help := &cobra.Command{Use: "help"}
	if !shouldSkipStartupInit(help) {
		t.Fatal("expected help command to skip startup init")
	}

	completion := &cobra.Command{Use: "completion"}
	if !shouldSkipStartupInit(completion) {
		t.Fatal("expected completion command to skip startup init")
	}

	collect := &cobra.Command{Use: "collect"}
	if shouldSkipStartupInit(collect) {
		t.Fatal("expected collect command to require startup init")
	}
}

func TestShouldSkipStartupInit_HelpFlag(t *testing.T) {
	collect := &cobra.Command{Use: "collect"}
	collect.Flags().Bool("help", false, "help for collect")
	if err := collect.Flags().Set("help", "true"); err != nil {
		t.Fatalf("failed to set help flag: %v", err)
	}

	if !shouldSkipStartupInit(collect) {
		t.Fatal("expected help flag to skip startup init")
	}
}

func TestInitializeForCommand_SkipsConfigForVersion(t *testing.T) {
	viper.Reset()
	cfgFile = ""
	strictProxy = false

	cmd := &cobra.Command{Use: "version"}
	if err := initializeForCommand(cmd); err != nil {
		t.Fatalf("initializeForCommand failed: %v", err)
	}

	if got := viper.ConfigFileUsed(); got != "" {
		t.Fatalf("expected no config file to be loaded for version, got %q", got)
	}
}
