package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func init() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel) // Set default log level to INFO

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		level, err := zerolog.ParseLevel(logLevel)
		if err == nil {
			zerolog.SetGlobalLevel(level)
		}
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Include code location if log level is DEBUG
	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		logger = logger.With().Caller().Logger()
	} else {
		logger = logger.With().Logger()
	}

	// Set the logger with the proper context to the global Logger
	Logger = logger
}
