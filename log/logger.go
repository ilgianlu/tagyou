package log

import (
	"log/slog"
	"os"
)

func Init() {
	minLevel := slog.LevelInfo
	if os.Getenv("DEBUG") != "" {
		minLevel = slog.LevelDebug
	}

	loggerOptions := slog.HandlerOptions{
		Level: minLevel,
	}
	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&loggerOptions,
		),
	)
	slog.SetDefault(logger)
}
