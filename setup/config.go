package setup

import (
	"fmt"
	"log/slog"

	viper "github.com/spf13/viper"
)

func SetupConfig(logger *slog.Logger) {
	viper.SetConfigName("base")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")

	// Setting prefix for all env variables: SCAFFOLD_ID => viper.Get("ID")
	viper.SetEnvPrefix("SCAFFOLD")

	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("Failed to load configurations", slog.String("side", "client"))
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// logger.Info("Configuration loaded successfully...", slog.String("side", "client"))
}
