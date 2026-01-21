package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// InitLogger initializes the global logger.
func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logLevel := viper.GetString("logging.level")
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	var output io.Writer
	logFormat := viper.GetString("logging.format")

	if logFormat == "text" {
		output = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	} else {
		// Ensure log directory exists if we were to log to file, 
		// but for now, let's stick to stdout/stderr based on config or file output if configured.
		// The requirement said "Ensure logs are written to both stderr (human-readable) and a file (JSON)."
		// Let's implement multi-writer.
		
		// Create logs directory
		if err := os.MkdirAll("logs", 0755); err != nil {
			panic(err)
		}
		
		file, err := os.OpenFile(filepath.Join("logs", "spectre.json"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			panic(err)
		}

		consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
		output = zerolog.MultiLevelWriter(consoleWriter, file)
	}
    
    // Fallback if not specifically "text" but "json" requested, or default
    if output == nil {
         output = os.Stderr
    }

	log.Logger = zerolog.New(output).With().Timestamp().Logger()
}
