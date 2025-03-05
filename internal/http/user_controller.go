package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/macreai/chess-game-app-be/internal/model"
	"github.com/macreai/chess-game-app-be/internal/usecase"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log     *logrus.Logger
	UseCase *usecase.UserUseCase
}

func NewUserController(log *logrus.Logger, usecase *usecase.UserUseCase) *UserController {
	return &UserController{
		Log:     log,
		UseCase: usecase,
	}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response := c.UseCase.Register(ctx.UserContext(), request)

	return ctx.JSON(response)
}
