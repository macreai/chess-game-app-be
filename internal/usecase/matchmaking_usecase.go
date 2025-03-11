package usecase

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/macreai/chess-game-app-be/internal/model"
	"github.com/macreai/chess-game-app-be/internal/repo"
	"github.com/sirupsen/logrus"
)

type MatchMakingUsecase struct {
	Log                   *logrus.Logger
	Validate              *validator.Validate
	MatchmakingRepository *repo.MatchMakingRepositoryImpl
}

func NewMatchMakingUsecase(
	log *logrus.Logger,
	validate *validator.Validate,
	matchMakingRepository *repo.MatchMakingRepositoryImpl,
) *MatchMakingUsecase {
	return &MatchMakingUsecase{
		Log:                   log,
		Validate:              validate,
		MatchmakingRepository: matchMakingRepository,
	}
}

func (m *MatchMakingUsecase) StartMatchMaking(request *model.CreateMatchMakingRequest) *model.WebResponse[*model.CreateMatchMakingResponse] {
	err := m.Validate.Struct(request)
	if err != nil {
		m.Log.Warnf("Invalid request body: %+v", err)
		return &model.WebResponse[*model.CreateMatchMakingResponse]{
			Errors: fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid request body: %+v", err)),
			Data:   nil,
			Status: fiber.StatusBadRequest,
		}
	}

	err = m.MatchmakingRepository.AddUserToQueue(request.UserID)
	if err != nil {
		m.Log.Warnf("failed to add user in queue: %+v", err)
		return &model.WebResponse[*model.CreateMatchMakingResponse]{
			Errors: fiber.NewError(fiber.StatusInternalServerError, "Failed to start matchmaking"),
			Data:   nil,
			Status: fiber.StatusInternalServerError,
		}
	}

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	var users []string

	for range ticker.C {
		count, err := m.MatchmakingRepository.MatchingUsers()
		if err != nil {
			m.Log.Errorf("failed to check q length: %+v", err)
			return &model.WebResponse[*model.CreateMatchMakingResponse]{
				Errors: fiber.NewError(fiber.StatusInternalServerError, "Error getting request"),
				Data:   nil,
				Status: fiber.StatusInternalServerError,
			}
		}

		if count >= 3 {
			users, err = m.MatchmakingRepository.MatchedUsers()
			if err != nil {
				m.Log.Errorf("failed to pop and get users: %+v", err)
				return &model.WebResponse[*model.CreateMatchMakingResponse]{
					Errors: fiber.NewError(fiber.StatusInternalServerError, "failed to get the users data"),
					Data:   nil,
					Status: fiber.StatusInternalServerError,
				}
			}

			break
		}
	}

	err = m.MatchmakingRepository.PopUsersMatched(users)
	if err != nil {
		m.Log.Warnf("failed to pop users matched: %+v", err)
	}

	roomId, err := uuid.NewUUID()
	if err != nil {
		m.Log.Errorf("failed to create room id: %+v", err)
		return &model.WebResponse[*model.CreateMatchMakingResponse]{
			Errors: fiber.NewError(fiber.StatusInternalServerError, "failed to create room id"),
			Data:   nil,
			Status: fiber.StatusInternalServerError,
		}
	}

	return &model.WebResponse[*model.CreateMatchMakingResponse]{
		Errors: nil,
		Data: &model.CreateMatchMakingResponse{
			RoomID: roomId.String(),
			Users:  users,
		},
		Status: fiber.StatusOK,
	}

}
