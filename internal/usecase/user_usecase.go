package usecase

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/macreai/chess-game-app-be/internal/entity"
	"github.com/macreai/chess-game-app-be/internal/model"
	"github.com/macreai/chess-game-app-be/internal/repo"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repo.UserRepositoryImpl
}

func NewUserUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate, userRepository *repo.UserRepositoryImpl) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            logger,
		Validate:       validate,
		UserRepository: userRepository,
	}
}

func (c *UserUseCase) Register(ctx context.Context, request *model.RegisterUserRequest) *model.WebResponse[*model.RegisterUserResponse] {
	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return &model.WebResponse[*model.RegisterUserResponse]{
			Errors: fiber.NewError(fiber.ErrBadRequest.Code, fmt.Sprintf("Invalid request body: %+v", err)),
			Data:   nil,
		}
	}

	_, err = c.UserRepository.FindByUsername(c.DB, request.Username)
	if err != gorm.ErrRecordNotFound {
		c.Log.Warnf("User : %+v", err)
		c.Log.Warnf("User already exists : %+v", err)
		c.Log.Tracef("Request Username : %v", request.Username)
		return &model.WebResponse[*model.RegisterUserResponse]{
			Errors: fiber.NewError(fiber.ErrConflict.Code, fmt.Sprintf("Username already exist!")),
			Data:   nil,
		}
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to generate bcrypt hash : %+v", err)
		return &model.WebResponse[*model.RegisterUserResponse]{
			Errors: fiber.ErrInternalServerError,
			Data:   nil,
		}
	}

	user := &entity.User{
		Name:     request.Name,
		Username: request.Username,
		Password: string(password),
	}

	if err := c.UserRepository.Create(c.DB, user); err != nil {
		c.Log.Warnf("Failed create user to database : %+v", err)
		return &model.WebResponse[*model.RegisterUserResponse]{
			Errors: fiber.ErrInternalServerError,
			Data:   nil,
		}
	}

	return &model.WebResponse[*model.RegisterUserResponse]{
		Errors: nil,
		Data: &model.RegisterUserResponse{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.Name,
		},
	}

}
