package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/macreai/chess-game-app-be/internal/auth"
	"github.com/macreai/chess-game-app-be/internal/entity"
	"github.com/macreai/chess-game-app-be/internal/model"
	"github.com/macreai/chess-game-app-be/internal/repo"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repo.UserRepositoryImpl
	Auth           *auth.MyJWT
	RedisDB        *redis.Client
}

func NewUserUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	validate *validator.Validate,
	userRepository *repo.UserRepositoryImpl,
	auth *auth.MyJWT,
	redisDb *redis.Client,
) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            logger,
		Validate:       validate,
		UserRepository: userRepository,
		Auth:           auth,
		RedisDB:        redisDb,
	}
}

func (c *UserUseCase) Register(request *model.RegisterUserRequest) *model.WebResponse[*model.RegisterUserResponse] {
	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return &model.WebResponse[*model.RegisterUserResponse]{
			Errors: fiber.NewError(fiber.ErrBadRequest.Code, fmt.Sprintf("Invalid request body: %+v", err)),
			Data:   nil,
			Status: fiber.StatusBadRequest,
		}
	}

	_, err = c.UserRepository.FindByUsername(c.DB, request.Username)
	if err != gorm.ErrRecordNotFound {
		c.Log.Warnf("Username already exist! %+v", err)
		return &model.WebResponse[*model.RegisterUserResponse]{
			Errors: fiber.NewError(fiber.ErrConflict.Code, fmt.Sprintf("Username already exist! %+v", err)),
			Data:   nil,
			Status: fiber.StatusConflict,
		}
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to generate bcrypt hash : %+v", err)
		return &model.WebResponse[*model.RegisterUserResponse]{
			Errors: fiber.ErrInternalServerError,
			Data:   nil,
			Status: fiber.StatusInternalServerError,
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
			Status: fiber.StatusInternalServerError,
		}
	}

	return &model.WebResponse[*model.RegisterUserResponse]{
		Errors: nil,
		Data: &model.RegisterUserResponse{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.Name,
		},
		Status: fiber.StatusOK,
	}

}

func (c *UserUseCase) Login(request *model.LoginUserRequest) *model.WebResponse[*model.LoginUserResponse] {
	ctxBackground := context.Background()
	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return &model.WebResponse[*model.LoginUserResponse]{
			Errors: fiber.NewError(fiber.ErrBadRequest.Code, fmt.Sprintf("Invalid request body: %+v", err)),
			Data:   nil,
			Status: fiber.StatusBadRequest,
		}
	}

	user, err := c.UserRepository.FindByUsername(c.DB, request.Usename)
	if err != nil {
		c.Log.Warnf("User not found: %+v", err)
		return &model.WebResponse[*model.LoginUserResponse]{
			Errors: fiber.NewError(fiber.ErrUnauthorized.Code, fmt.Sprintf("User not found: %+v", err)),
			Data:   nil,
			Status: fiber.StatusUnauthorized,
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return &model.WebResponse[*model.LoginUserResponse]{
			Errors: fiber.NewError(fiber.ErrUnauthorized.Code, fmt.Sprintf("Invalid credentias: %+v", err)),
			Data:   nil,
			Status: fiber.StatusUnauthorized,
		}
	}

	exp, token, err := c.Auth.GenerateJWT(user, c.Auth.Viper)

	if err != nil {
		c.Log.Warnf("Error generate JWT: %+v", err)
		return &model.WebResponse[*model.LoginUserResponse]{
			Errors: fiber.NewError(fiber.ErrInternalServerError.Code, fmt.Sprintf("Error generate JWT: %+v", err)),
			Data:   nil,
			Status: fiber.StatusUnauthorized,
		}
	}

	c.RedisDB.Set(ctxBackground, fmt.Sprintf("%d", user.ID), token, time.Until(exp))

	return &model.WebResponse[*model.LoginUserResponse]{
		Errors: nil,
		Data: &model.LoginUserResponse{
			Token: token,
		},
		Status: fiber.StatusOK,
	}
}
