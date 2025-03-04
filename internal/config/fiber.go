package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func NewFiber(viper *viper.Viper) *fiber.App {
	fiber := fiber.New(fiber.Config{
		AppName:      viper.GetString("APP_NAME"),
		Prefork:      viper.GetBool("WEB_PREFORK"),
		ErrorHandler: NewErrorHandler(),
	})

	return fiber
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		return ctx.Status(code).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}
}
