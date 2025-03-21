package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/macreai/chess-game-app-be/internal/auth"
	"github.com/macreai/chess-game-app-be/internal/http"
	"github.com/macreai/chess-game-app-be/internal/http/route"
	"github.com/macreai/chess-game-app-be/internal/repo"
	"github.com/macreai/chess-game-app-be/internal/usecase"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type AppConfig struct {
	App       *fiber.App
	DB        *gorm.DB
	Log       *logrus.Logger
	Validator *validator.Validate
	Config    *viper.Viper
	Jwt       *auth.MyJWT
	RedisDB   *redis.Client
	Viper     *viper.Viper
}

func InitApp(app *AppConfig) {
	userRepository := repo.NewUserRepositoryImpl(app.Log)
	userUsecase := usecase.NewUserUseCase(app.DB, app.Log, app.Validator, userRepository, app.Jwt)
	userController := http.NewUserController(app.Log, userUsecase)

	matchRepository := repo.NewMatchMakingRepositoryImpl(app.Log, app.RedisDB, app.Viper)
	matchUsecase := usecase.NewMatchMakingUsecase(app.Log, app.Validator, matchRepository)
	matchController := http.NewMatchMakingController(app.Log, matchUsecase)

	routeConfig := &route.RouteConfig{
		App:                   app.App,
		UserController:        userController,
		MatchMakingController: matchController,
		AuthMiddleware:        app.Jwt.JWTMiddleware(app.Config),
	}
	routeConfig.Setup()
}
