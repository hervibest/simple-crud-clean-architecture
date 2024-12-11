package entity

import (
	"time"

	"github.com/google/uuid"
)

type CourseCategory struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement"`
	UUID        uuid.UUID `gorm:"column:uuid"`
	Name        string    `gorm:"column:name"`
	Slug        string    `gorm:"column:slug"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Course      []Course  `gorm:"many2many:course_category_course;foreignKey:id;joinForeignKey:course_category_id;references:id;joinReferences:course_id"`
}

func (u *CourseCategory) TableName() string {
	return "course_categories"
}
