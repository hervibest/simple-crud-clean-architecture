package entity

import (
	"time"

	"github.com/google/uuid"
)

type SectionVideo struct {
	ID           int       `gorm:"column:id;primaryKey;autoIncrement"`
	UUID         uuid.UUID `gorm:"column:uuid"`
	SectionID    int       `gorm:"column:section_id"`
	Title        string    `gorm:"column:title"`
	Sequence     int       `gorm:"column:sequence"`
	Notes        string    `gorm:"column:notes"`
	OriginalName string    `gorm:"column:original_name"`
	OriginalSize float64   `gorm:"column:original_size"`
	OriginalMime string    `gorm:"column:original_mime"`
	MediaID      string    `gorm:"column:media_id"`
	Bucket       string    `gorm:"column:bucket"`
	Dir          string    `gorm:"column:dir"`
	Video        *File     `gorm:"polymorphic:Fileable;polymorphicValue:Video"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (u *SectionVideo) TableName() string {
	return "section_videos"
}
