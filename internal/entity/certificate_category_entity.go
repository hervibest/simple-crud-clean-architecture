package entity

import (
	"time"

	"github.com/google/uuid"
)

type CertificateCategory struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement"`
	UUID        uuid.UUID `gorm:"column:uuid"`
	Name        string    `gorm:"column:name"`
	Slug        string    `gorm:"column:slug"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (u *CertificateCategory) TableName() string {
	return "certificate_categories"
}
