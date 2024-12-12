package entity

import (
	"time"

	"github.com/google/uuid"
)

type CourseSection struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement"`
	UUID        uuid.UUID `gorm:"column:uuid"`
	CourseID    int       `gorm:"column:course_id"`
	Title       string    `gorm:"column:title"`
	Sequence    int       `gorm:"column:sequence"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (u *CourseSection) TableName() string {
	return "course_sections"
}
