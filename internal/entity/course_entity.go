package entity

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID          string           `gorm:"column:id;primaryKey;autoIncrement"`
	UUID        uuid.UUID        `gorm:"column:uuid"`
	Name        string           `gorm:"column:name"`
	Slug        string           `gorm:"column:slug"`
	Description string           `gorm:"column:description"`
	Price       float64          `gorm:"column:price"`
	IsActive    bool             `gorm:"column:is_active"`
	CreatedAt   time.Time        `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time        `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Categories  []CourseCategory `gorm:"many2many:course_category_course;foreignKey:id;joinForeignKey:course_id;references:id;joinReferences:course_category_id"`
}

func (u *Course) TableName() string {
	return "courses"
}
