package test

import (
	"testing"

	"github.com/macreai/chess-game-app-be/internal/config"
)

func TestInit(t *testing.T) {
	viperConfig := config.NewViper()
	log := config.NewLogrus(viperConfig)
	validate := config.NewValidator()
	app := config.NewFiber(viperConfig)
	db := config.NewDatabase(viperConfig, log)

	config.InitApp(&config.AppConfig{
		App:       app,
		DB:        db,
		Log:       log,
		Validator: validate,
		Config:    viperConfig,
	})
}
