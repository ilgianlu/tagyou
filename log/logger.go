package log

import (
	"log/slog"
	"os"
)

func Init(debugging bool) slog.Level {
	minLevel := slog.LevelInfo
	if debugging {
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
	slog.SetLogLoggerLevel(minLevel)
	slog.SetDefault(logger)
	return minLevel
}
