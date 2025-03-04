package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewLogrus(viper *viper.Viper) *logrus.Logger {
	logger := logrus.New()

	logger.SetLevel(logrus.Level(viper.GetInt("LOG_LEVEL")))
	logger.SetFormatter(&logrus.JSONFormatter{})

	return logger
}
