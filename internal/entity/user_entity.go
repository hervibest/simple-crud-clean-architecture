package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           int           `gorm:"column:id;primaryKey;autoIncrement"`
	UUID         uuid.UUID     `gorm:"column:uuid"`
	Password     string        `gorm:"column:password"`
	Email        string        `gorm:"column:email"`
	CreatedAt    time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time     `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	VerifiedAt   *time.Time    `gorm:"column:verified_at"`
	AccessToken  string        `gorm:"-"`
	RefreshToken string        `gorm:"-"`
	Transaction  []Transaction `gorm:"foreignKey:user_id;references:id"`
	Courses      []Course      `gorm:"many2many:course_user;foreignKey:id;joinForeignKey:user_id;references:id;joinReferences:course_id"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) HasVerifiedEmail() bool {
	return u.VerifiedAt != nil
}
