package entity

import (
	"time"

	"github.com/google/uuid"
)

type Material struct {
	ID            int         `gorm:"column:id;primaryKey;autoIncrement"`
	UUID          uuid.UUID   `gorm:"column:uuid"`
	CertificateID int         `gorm:"column:certificate_id"`
	Name          string      `gorm:"column:name"`
	Code          string      `gorm:"column:code"`
	Type          string      `gorm:"column:type"`
	CreatedAt     time.Time   `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time   `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Certificates  Certificate `gorm:"foreignKey:certificate_id;references:id"`
	Thumbnail     *File       `gorm:"polymorphic:Fileable;polymorphicValue:MaterialThumbnail"`
}

func (u *Material) TableName() string {
	return "certificates"
}
