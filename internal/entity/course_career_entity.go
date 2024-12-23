package entity

import (
	"time"

	"github.com/google/uuid"
)

type CareerCategory struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement"`
	UUID        uuid.UUID `gorm:"column:uuid"`
	Name        string    `gorm:"column:name"`
	Slug        string    `gorm:"column:slug"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Career      []Career  `gorm:"many2many:career_category_career;foreignKey:id;joinForeignKey:career_category_id;references:id;joinReferences:career_id"`
}

func (u *CareerCategory) TableName() string {
	return "career_categories"
}
