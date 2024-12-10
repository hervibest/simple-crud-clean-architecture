package entity

import "time"

type VerificationToken struct {
	Token     string    `gorm:"column:token"`
	UserEmail string    `gorm:"column:email;primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	User      User      `gorm:"foreignKey:UserEmail;references:Email"`
}

func (c *VerificationToken) TableName() string {
	return "user_verification_tokens"
}
