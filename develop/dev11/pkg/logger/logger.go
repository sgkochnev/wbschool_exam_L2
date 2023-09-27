package logger

import (
	"log/slog"
	"os"
)

// New - create logger
func New(logLevel slog.Level) *slog.Logger {
	var log *slog.Logger
	log = slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: logLevel,
			},
		),
	)
	return log
}
