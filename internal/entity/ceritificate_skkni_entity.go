package entity

import (
	"time"

	"github.com/google/uuid"
)

type Skkni struct {
	ID            int         `gorm:"column:id;primaryKey;autoIncrement"`
	UUID          uuid.UUID   `gorm:"column:uuid"`
	CertificateID int         `gorm:"column:certificate_id"`
	Name          string      `gorm:"column:name"`
	CreatedAt     time.Time   `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time   `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Certificates  Certificate `gorm:"foreignKey:certificate_id;references:id"`
	File          *File       `gorm:"polymorphic:Fileable;polymorphicValue:SkkniFile"`
}

func (u *Skkni) TableName() string {
	return "skkni"
}
