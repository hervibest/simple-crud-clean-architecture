package model

import (
	"time"

	"github.com/google/uuid"
)

type SearchCertificateMatRequest struct {
	CertificateID   int
	CertificateUUID uuid.UUID `json:"certificate_uuid" validate:"required"`
	Name            string
	Code            string
	Type            string
	Page            int
	Size            int
}

type GetCertificateMatRequest struct {
	UUID uuid.UUID `json:"uuid,omitempty" validate:"required"`
}

type CreateCertificateMatRequest struct {
	CertificateUUID uuid.UUID `json:"certificate_uuid" validate:"required"`
	Name            string    `json:"name" validate:"required,max=255"`
	Code            string    `json:"code" validate:"required,max=255"`
	Type            string    `json:"type" validate:"required,max=255"`
	Description     string    `json:"description"`
}

type UpdateCertificateMatRequest struct {
	CertificateUUID    uuid.UUID `json:"certificate_uuid" validate:"required"`
	CertificateMatUUID uuid.UUID `validate:"required"`
	Name               string    `json:"name" validate:"required,max=255"`
	Code               string    `json:"code"`
	Type               string    `json:"type"`
}

type DeleteCertificateMatRequest struct {
	CertificateUUID    uuid.UUID `json:"certificate_uuid" validate:"required"`
	CertificateMatUUID uuid.UUID `validate:"required"`
}
type CertificateMaterialResponse struct {
	UUID      uuid.UUID `json:"uuid,omitempty"`
	Name      string    `json:"name,omitempty"`
	Code      string    `json:"code,omitempty"`
	Type      string    `json:"type,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
