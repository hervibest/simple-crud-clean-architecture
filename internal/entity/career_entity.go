package entity

import (
	"time"

	"github.com/google/uuid"
)

type Career struct {
	ID          int              `gorm:"column:id;primaryKey;autoIncrement"`
	UUID        uuid.UUID        `gorm:"column:uuid"`
	Title       string           `gorm:"column:title"`
	Slug        string           `gorm:"column:slug"`
	Description string           `gorm:"column:description"`
	Price       float64          `gorm:"column:price"`
	FinalPrice  float64          `gorm:"-"`
	IsActive    bool             `gorm:"column:is_active"`
	CreatedAt   time.Time        `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time        `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Categories  []CareerCategory `gorm:"many2many:career_category_career;foreignKey:id;joinForeignKey:career_id;references:id;joinReferences:career_category_id"`
	Thumbnail   *File            `gorm:"polymorphic:Fileable;polymorphicValue:CareerThumbnail"`
}

func (u *Career) TableName() string {
	return "careers"
}
