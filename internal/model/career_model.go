package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateCareerRequest struct {
	Title         string `json:"title" validate:"required,max=255"`
	Slug          string
	Description   string   `json:"description"`
	Price         float64  `json:"price"`
	Categories    []string `json:"categories"`
	CategoryUUIDs []uuid.UUID
}

type GetCareerRequest struct {
	UUID uuid.UUID `json:"uuid,omitempty" validate:"required"`
}

type GetPurchasedCareerRequest struct {
	UserID int `validate:"required"`
}

type UpdateCareerRequest struct {
	UUID          uuid.UUID `json:"uuid,omitempty" validate:"required"`
	Title         string    `json:"title" validate:"required,max=255"`
	Slug          string
	Description   string   `json:"description"`
	Price         float64  `json:"price"`
	IsActive      bool     `json:"is_active"`
	Categories    []string `json:"categories"`
	CategoryUUIDs []uuid.UUID
}

type CareerResponse struct {
	UUID         uuid.UUID           `json:"uuid,omitempty"`
	Title        string              `json:"title,omitempty"`
	Slug         string              `json:"slug,omitempty"`
	Description  string              `json:"description,omitempty"`
	Price        float64             `json:"original_price,omitempty"`
	FinalPrice   float64             `json:"final_price,omitempty"`
	IsActive     bool                `json:"is_active,omitempty"`
	Categories   []CareerCatResponse `json:"categories,omitempty"`
	ThumbnailURL string              `json:"thumbnail_url,omitempty"`
	CreatedAt    time.Time           `json:"created_at,omitempty"`
	UpdatedAt    time.Time           `json:"updated_at,omitempty"`
}

type SearchCareerRequest struct {
	Title       string `json:"title" validate:"max=255"`
	Slug        string `json:"slug" validate:"max=255"`
	Description string `json:"description" validate:"max=255"`
	Page        int    `json:"page" validate:"min=1"`
	Size        int    `json:"size" validate:"min=1,max=100"`
}

type SearchPurchasedCareer struct {
	UserID      int    `validate:"required"`
	Title       string `json:"title" validate:"max=255"`
	Slug        string `json:"slug" validate:"max=255"`
	Description string `json:"description" validate:"max=255"`
	Page        int    `json:"page" validate:"min=1"`
	Size        int    `json:"size" validate:"min=1,max=100"`
}

type CareerUUIDresponse struct {
	UUID uuid.UUID `json:"uuid,omitempty"`
}

type CareerThumbnailRequest struct {
	CareerUUID uuid.UUID `validate:"required"`
}
