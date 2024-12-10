package entity

import (
	"time"

	"github.com/google/uuid"
)

// User is a struct that represents a user entity
type User struct {
	ID         string     `gorm:"column:id;primaryKey;autoIncrement"`
	UUID       uuid.UUID  `gorm:"column:uuid"`
	Password   string     `gorm:"column:password"`
	Email      string     `gorm:"column:email"`
	CreatedAt  time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time  `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	VerifiedAt *time.Time `gorm:"column:verified_at"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) HasVerifiedEmail() bool {
	return u.VerifiedAt != nil
}
