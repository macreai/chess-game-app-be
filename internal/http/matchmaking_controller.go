package http

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/macreai/chess-game-app-be/internal/model"
	"github.com/macreai/chess-game-app-be/internal/usecase"
	"github.com/sirupsen/logrus"
)

type MatchMakingController struct {
	Log                *logrus.Logger
	MatchMakingUsecase *usecase.MatchMakingUsecase
}

func NewMatchMakingController(log *logrus.Logger, matchMakingUsecase *usecase.MatchMakingUsecase) *MatchMakingController {
	return &MatchMakingController{
		Log:                log,
		MatchMakingUsecase: matchMakingUsecase,
	}
}

func (m *MatchMakingController) StartMatchMaking(ctx *fiber.Ctx) error {
	request := new(model.CreateMatchMakingRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		m.Log.Warnf("Failed to parse request body : %+v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[*model.CreateMatchMakingResponse]{
			Errors: fiber.NewError(fiber.ErrBadRequest.Code, fmt.Sprintf("Failed to parse request body: %+v", err)),
			Data:   nil,
			Status: fiber.StatusBadRequest,
		})
	}

	response := m.MatchMakingUsecase.StartMatchMaking(request)

	return ctx.Status(fiber.StatusOK).JSON(response)
}
