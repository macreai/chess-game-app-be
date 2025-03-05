package entity

import "time"

type User struct {
	ID        uint64    `gorm:"column:id;primaryKey"`
	Username  string    `gorm:"column:username"`
	Password  string    `gorm:"column:password"`
	Name      string    `gorm:"column:name"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Scores    []Score   `gorm:"foreignKey:UserID"`
}

func (u *User) TableName() string {
	return "users"
}
