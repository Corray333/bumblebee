package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Corray333/bumblebee/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func MustInit(envPath string) {
	if err := godotenv.Load(envPath); err != nil {
		panic("error while loading .env file: " + err.Error())
	}
	SetupLogger()
	configPath := fmt.Sprintf("../configs/%s.yaml", os.Getenv("ENV"))
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func SetupLogger() {
	handler := logger.NewHandler(nil)
	log := slog.New(handler)
	slog.SetDefault(log)
}
