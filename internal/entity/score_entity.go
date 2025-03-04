package entity

import "time"

type Score struct {
	ID    uint64    `gorm:"primaryKey"`
	User  User      `gorm:"foreignKey:user_id;references:id"`
	Score int       `gorm:"column:user"`
	Date  time.Time `gorm:"column:date;autoCreateTime"`
}

func (s *Score) TableName() string {
	return "scores"
}
