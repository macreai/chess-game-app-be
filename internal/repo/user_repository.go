package repo

import (
	"github.com/macreai/chess-game-app-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository interface {
	Register(db *gorm.DB, user *entity.User)
	GetUser(db *gorm.DB, user *entity.User)
}

type UserRepositoryImpl struct {
	Repository[entity.User]
	Log *logrus.Logger
}

func NewUserRepositoryImpl(log *logrus.Logger) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		Log: log,
	}
}

func (u *UserRepositoryImpl) Register(db *gorm.DB, user *entity.User) error {
	return u.Repository.Create(db, user)
}

func (u *UserRepositoryImpl) GetUser(db *gorm.DB, user *entity.User) error {
	return u.Repository.FindById(db, user, user.ID)
}
