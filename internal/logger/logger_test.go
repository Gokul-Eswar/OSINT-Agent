package logger

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestInitLogger(t *testing.T) {
	// Set some config for the logger
	viper.Set("logging.level", "debug")
	viper.Set("logging.format", "text")

	// Initialize logger (should not panic)
	InitLogger()
	
	// Clean up logs directory if created
	defer os.RemoveAll("logs")
}
