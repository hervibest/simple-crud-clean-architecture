package entity

import (
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement"`
	UUID        uuid.UUID `gorm:"column:uuid"`
	Password    string    `gorm:"column:password"`
	Email       string    `gorm:"column:email"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	AccessToken string    `gorm:"-"`
	// RefreshToken string    `gorm:"-"`
}

func (u *Employee) TableName() string {
	return "employees"
}
