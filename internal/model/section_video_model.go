package model

import (
	"time"

	"github.com/google/uuid"
)

type SearchSecVideosRequest struct {
	SectionID   int
	SectionUUID uuid.UUID `json:"section_uuid" validate:"required"`
	Title       string    `json:"title" `
	Notes       string    `json:"notes" `
	Page        int
	Size        int
}

type GetSecVideoRequest struct {
	UUID uuid.UUID `json:"uuid,omitempty" validate:"required"`
}

type CreateSecVideoRequest struct {
	SectionUUID uuid.UUID `json:"section_uuid" validate:"required"`
	Title       string    `json:"title" validate:"required,max=255"`
	Notes       string    `json:"notes" validate:"required"`
	Sequence    int       `json:"sequence" validate:"required"`
	Description string    `json:"description"`
}

type UpdateSecVideoRequest struct {
	SectionUUID uuid.UUID `json:"section_uuid" validate:"required"`
	VideoUUID   uuid.UUID `validate:"required"`
	Title       string    `json:"title" validate:"required,max=255"`
	Notes       string    `json:"notes" validate:"required"`
	Sequence    int       `json:"sequence" validate:"required"`
	Description string    `json:"description"`
}

type DeleteSecVideoRequest struct {
	SectionUUID uuid.UUID `json:"section_uuid" validate:"required"`
	VideoUUID   uuid.UUID `validate:"required"`
}
type SecVideoResponse struct {
	UUID         uuid.UUID `json:"uuid,omitempty"`
	Title        string    `json:"title,omitempty"`
	Notes        string    `json:"notes,omitempty"`
	Sequence     int       `json:"sequence,omitempty"`
	OriginalName string    `json:"original_name,omitempty"`
	OriginalSize float64   `json:"original_size,omitempty"`
	OriginalMime string    `json:"original_mimeomitempty"`
	MediaID      string    `json:"media_id,omitempty"`
	Bucket       string    `json:"bucket,omitempty"`
	Dir          string    `json:"dir,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}
