package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/macreai/chess-game-app-be/internal/http"
)

type RouteConfig struct {
	App            *fiber.App
	UserController *http.UserController
}

func (r *RouteConfig) Setup() {
	r.setupGuestRoute()
}

func (r *RouteConfig) setupGuestRoute() {
	r.App.Post("/api/v1/users/register", r.UserController.Register)
}
