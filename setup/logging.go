package setup

import (
	"log/slog"
	"os"
)

func SetupLogging() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// logger.Info("Initialized logging", slog.String("side", "client"))

	return logger
}
