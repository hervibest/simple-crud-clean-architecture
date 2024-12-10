package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateCourseCatRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Slug        string
	Description string `json:"description"`
}

type GetCourseCatRequest struct {
	UUID uuid.UUID `json:"uuid,omitempty" validate:"required"`
}

type UpdateCourseCatRequest struct {
	UUID        uuid.UUID `json:"uuid,omitempty" validate:"required"`
	Name        string    `json:"name" validate:"required,max=255"`
	Slug        string
	Description string `json:"description"`
}

type CourseCatResponse struct {
	UUID        uuid.UUID `json:"uuid,omitempty"`
	Name        string    `json:"email,omitempty"`
	Slug        string    `json:"slug,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
