package entity

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID           int       `gorm:"column:id;primaryKey;autoIncrement"`
	UUID         uuid.UUID `gorm:"column:uuid"`
	Filename     string    `gorm:"column:filename"`
	Mimetype     string    `gorm:"column:mimetype"`
	Path         string    `gorm:"column:path"`
	Size         int64     `gorm:"column:size"`
	FileableID   int       `gorm:"column:fileable_id"`
	FileableType string    `gorm:"column:fileable_type" `
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
