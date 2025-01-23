package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateCertificateRequest struct {
	Name         string `json:"name" validate:"required,max=255"`
	Slug         string
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	Category     string  `json:"category"`
	CategoryUUID uuid.UUID
}

type GetCertificateRequest struct {
	UUID uuid.UUID `json:"uuid,omitempty" validate:"required"`
}

type GetPurchasedCertificateRequest struct {
	UserID int `validate:"required"`
}

type UpdateCertificateRequest struct {
	UUID         uuid.UUID `json:"uuid,omitempty" validate:"required"`
	Name         string    `json:"name" validate:"required,max=255"`
	Slug         string
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	IsActive     bool    `json:"is_active"`
	Category     string  `json:"category"`
	CategoryUUID uuid.UUID
}

type CertificateResponse struct {
	UUID         uuid.UUID               `json:"uuid,omitempty"`
	Name         string                  `json:"name,omitempty"`
	Slug         string                  `json:"slug,omitempty"`
	Description  string                  `json:"descriptions,omitempty"`
	Price        float64                 `json:"original_price,omitempty"`
	FinalPrice   float64                 `json:"final_price,omitempty"`
	IsActive     bool                    `json:"is_active,omitempty"`
	Category     *CertificateCatResponse `json:"categories,omitempty"`
	ThumbnailURL string                  `json:"thumbnail_url,omitempty"`
	Discount     *DiscountResponse       `json:"discount,omitempty"`
	CreatedAt    time.Time               `json:"created_at,omitempty"`
	UpdatedAt    time.Time               `json:"updated_at,omitempty"`
}

type SearchCertificateRequest struct {
	Name        string `json:"name" validate:"max=255"`
	Slug        string `json:"slug" validate:"max=255"`
	Description string `json:"description" validate:"max=255"`
	Page        int    `json:"page" validate:"min=1"`
	Size        int    `json:"size" validate:"min=1,max=100"`
}

type SearchPurchasedCertificate struct {
	UserID      int    `validate:"required"`
	Name        string `json:"name" validate:"max=255"`
	Slug        string `json:"slug" validate:"max=255"`
	Description string `json:"description" validate:"max=255"`
	Page        int    `json:"page" validate:"min=1"`
	Size        int    `json:"size" validate:"min=1,max=100"`
}

type CertificateUUIDresponse struct {
	UUID uuid.UUID `json:"uuid,omitempty"`
}

type CertificateThumbnailRequest struct {
	CertificateUUID uuid.UUID `validate:"required"`
}
