package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateCertifCatRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Slug        string
	Description string `json:"description"`
}

type GetCertificateCatRequest struct {
	UUID uuid.UUID `json:"uuid,omitempty" validate:"required"`
}

type UpdateCertificateCatRequest struct {
	UUID        uuid.UUID `json:"uuid,omitempty" validate:"required"`
	Name        string    `json:"name" validate:"required,max=255"`
	Slug        string
	Description string `json:"description"`
}

type CertificateCatResponse struct {
	UUID        uuid.UUID `json:"uuid,omitempty"`
	Name        string    `json:"name,omitempty"`
	Slug        string    `json:"slug,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type SearchCertificateCatRequest struct {
	Name        string `json:"name" validate:"max=255"`
	Slug        string `json:"slug" validate:"max=255"`
	Description string `json:"description" validate:"max=255"`
	Page        int    `json:"page" validate:"min=1"`
	Size        int    `json:"size" validate:"min=1,max=100"`
}
