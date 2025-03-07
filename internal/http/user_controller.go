package http

import (
	"fmt"

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
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[*model.LoginUserResponse]{
			Errors: fiber.NewError(fiber.ErrBadRequest.Code, fmt.Sprintf("Failed to parse request body: %+v", err)),
			Data:   nil,
		})
	}

	response := c.UseCase.Register(request)

	return ctx.Status(response.Status).JSON(response)
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[*model.LoginUserResponse]{
			Errors: fiber.NewError(fiber.ErrBadRequest.Code, fmt.Sprintf("Failed to parse request body: %+v", err)),
			Data:   nil,
		})
	}

	response := c.UseCase.Login(request)

	return ctx.Status(response.Status).JSON(response)
}
