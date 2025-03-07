package main

import (
	"fmt"

	"github.com/macreai/chess-game-app-be/internal/auth"
	"github.com/macreai/chess-game-app-be/internal/config"
)

func main() {
	viperConfig := config.NewViper("./")
	fiberApp := config.NewFiber(viperConfig)
	logConfig := config.NewLogrus(viperConfig)
	gormConfig := config.NewDatabase(viperConfig, logConfig)
	validatorConfig := config.NewValidator()
	authConfig := auth.NewMyJWT(viperConfig)
	redisConfig := config.NewRedis(viperConfig)

	config.InitApp(&config.AppConfig{
		App:       fiberApp,
		DB:        gormConfig,
		Log:       logConfig,
		Validator: validatorConfig,
		Config:    viperConfig,
		Jwt:       authConfig,
		RedisDB:   redisConfig,
	})

	webPort := viperConfig.GetInt("WEB_PORT")
	err := fiberApp.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		logConfig.Fatalf("Failed to start the server : %v", err)
	}
}
