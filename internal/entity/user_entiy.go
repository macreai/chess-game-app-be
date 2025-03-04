package entity

type User struct {
	ID        uint64  `gorm:"column:id;primaryKey"`
	Username  string  `gorm:"column:username"`
	Password  string  `gorm:"column:password"`
	Name      string  `gorm:"column:name"`
	CreatedAt int64   `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt int64   `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Scores    []Score `gorm:"foreignKey:user_id;references:id"`
}

func (u *User) TableName() string {
	return "users"
}
