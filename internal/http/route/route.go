package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/macreai/chess-game-app-be/internal/http"
)

type RouteConfig struct {
	App            *fiber.App
	UserController *http.UserController
	AuthMiddleware fiber.Handler
}

func (r *RouteConfig) Setup() {
	r.setupGuestRoute()
	r.setupAuthRoute()
}

func (r *RouteConfig) setupGuestRoute() {
	r.App.Post("/api/v1/users", r.UserController.Register)
	r.App.Post("/api/v1/users/login", r.UserController.Login)
}

func (r *RouteConfig) setupAuthRoute() {
	r.App.Use(r.AuthMiddleware)
	r.App.Delete("/api/v1/users", func(ctx *fiber.Ctx) error {
		return ctx.SendString("I am deleted")
	})
}
