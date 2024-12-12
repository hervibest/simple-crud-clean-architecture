package entity

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID             int              `gorm:"column:id;primaryKey;autoIncrement"`
	UUID           uuid.UUID        `gorm:"column:uuid"`
	Name           string           `gorm:"column:name"`
	Slug           string           `gorm:"column:slug"`
	Description    string           `gorm:"column:description"`
	Price          float64          `gorm:"column:price"`
	IsActive       bool             `gorm:"column:is_active"`
	CreatedAt      time.Time        `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time        `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Categories     []CourseCategory `gorm:"many2many:course_category_course;foreignKey:id;joinForeignKey:course_id;references:id;joinReferences:course_category_id"`
	Transaction    []Transaction    `gorm:"foreignKey:course_id;references:id"`
	CourseSections []CourseSection  `gorm:"foreignKey:course_id;references:id"`
	User           []User           `gorm:"many2many:course_user;foreignKey:id;joinForeignKey:course_id;references:id;joinReferences:user_id"`
}

func (u *Course) TableName() string {
	return "courses"
}
