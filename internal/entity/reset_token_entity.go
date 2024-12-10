package entity

import "time"

type ResetPasswordToken struct {
	Token     string    `gorm:"column:token"`
	UserEmail string    `gorm:"column:email;primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	User      User      `gorm:"foreignKey:UserEmail;references:Email"`
}

func (c *ResetPasswordToken) TableName() string {
	return "reset_password_tokens"
}
