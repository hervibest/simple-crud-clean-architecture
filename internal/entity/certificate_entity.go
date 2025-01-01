package entity

import (
	"time"

	"github.com/google/uuid"
)

type Certificate struct {
	ID           int                  `gorm:"column:id;primaryKey;autoIncrement"`
	UUID         uuid.UUID            `gorm:"column:uuid"`
	Name         string               `gorm:"column:name"`
	Slug         string               `gorm:"column:slug"`
	Description  string               `gorm:"column:description"`
	Price        float64              `gorm:"column:price"`
	FinalPrice   float64              `gorm:"-"`
	IsActive     bool                 `gorm:"column:is_active"`
	CategoriesId int                  `gorm:"column:categories_id"`
	CreatedAt    time.Time            `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time            `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Category     *CertificateCategory `gorm:"foreignKey:categories_id;references:id"`
	Materials    []Material           `gorm:"foreignKey:certificate_id;references:id"`
	Skkni        []Skkni              `gorm:"foreignKey:certificate_id;references:id"`
	Thumbnail    *File                `gorm:"polymorphic:Fileable;polymorphicValue:CertificateThumbnail"`
}

func (u *Certificate) TableName() string {
	return "certificates"
}
