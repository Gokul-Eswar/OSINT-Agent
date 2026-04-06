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

// InitLogger initializes the global logger using the full path.
// Kept for backward compatibility with existing call sites.
func InitLogger() {
	InitLoggerFull()
}

// InitLoggerLight initializes a low-overhead stderr logger.
func InitLoggerLight() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(parseLogLevel())
	log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
}

// InitLoggerFull initializes logger outputs based on config.
func InitLoggerFull() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(parseLogLevel())
	log.Logger = zerolog.New(buildFullOutput()).With().Timestamp().Logger()
}

func parseLogLevel() zerolog.Level {
	logLevel := viper.GetString("logging.level")
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	return level
}

func buildFullOutput() io.Writer {
	logFormat := viper.GetString("logging.format")

	if logFormat == "text" {
		return zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	}

	if err := os.MkdirAll("logs", 0755); err != nil {
		return os.Stderr
	}

	file, err := os.OpenFile(filepath.Join("logs", "spectre.json"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return os.Stderr
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	return zerolog.MultiLevelWriter(consoleWriter, file)
}
