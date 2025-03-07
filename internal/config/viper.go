package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewViper(configPath string) *viper.Viper {
	viper := viper.New()

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return viper
}
