package repo

import (
	"github.com/macreai/chess-game-app-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository interface {
	Register(db *gorm.DB, user *entity.User) error
	GetUser(db *gorm.DB, id uint64) (*entity.User, error)
	FindByUsername(db *gorm.DB, user string) (*entity.User, error)
}

type UserRepositoryImpl struct {
	Repository[entity.User]
	log *logrus.Logger
}

func NewUserRepositoryImpl(log *logrus.Logger) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		log: log,
	}
}

func (u *UserRepositoryImpl) Register(db *gorm.DB, user *entity.User) error {
	return u.Repository.Create(db, user)
}

func (u *UserRepositoryImpl) GetUser(db *gorm.DB, id uint64) (*entity.User, error) {
	return u.Repository.FindById(db, id)
}

func (u *UserRepositoryImpl) FindByUsername(db *gorm.DB, username string) (entity.User, int64) {
	var user entity.User
	rows := db.Where("username = ?", username).First(&user).RowsAffected
	return user, rows
}
