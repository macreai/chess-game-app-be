package entity

import "time"

type Score struct {
	ID     uint64    `gorm:"primaryKey"`
	UserID uint64    `gorm:"column:user_id"`
	User   User      `gorm:"foreignKey:UserID;references:id"`
	Score  int       `gorm:"column:user"`
	Date   time.Time `gorm:"column:date;autoCreateTime"`
}

func (s *Score) TableName() string {
	return "scores"
}
