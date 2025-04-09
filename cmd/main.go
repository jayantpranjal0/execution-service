package main

import (
	"go.uber.org/zap"
	"github.com/spf13/viper"
)

func main() {
	// Load configuration
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./config")

	if err:= viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any
	logger.Info("Logger initialized")

	// Start coordinator

}