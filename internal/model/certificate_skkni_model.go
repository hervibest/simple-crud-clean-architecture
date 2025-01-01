package model

import (
	"time"

	"github.com/google/uuid"
)

type SearchCertificateSkkniRequest struct {
	CertificateID   int
	CertificateUUID uuid.UUID `json:"certificate_uuid" validate:"required"`
	Name            string
	Page            int
	Size            int
}

type GetCertificateSkkniRequest struct {
	UUID uuid.UUID `json:"uuid,omitempty" validate:"required"`
}

type CreateCertificateSkkniRequest struct {
	CertificateUUID uuid.UUID `json:"certificate_uuid" validate:"required"`
	Name            string    `json:"name" validate:"required,max=255"`
}

type UpdateCertificateSkkniRequest struct {
	CertificateUUID      uuid.UUID `json:"certificate_uuid" validate:"required"`
	CertificateSkkniUUID uuid.UUID `validate:"required"`
	Name                 string    `json:"name" validate:"required,max=255"`
}

type DeleteCertificateSkkniRequest struct {
	CertificateUUID      uuid.UUID `json:"certificate_uuid" validate:"required"`
	CertificateSkkniUUID uuid.UUID `validate:"required"`
}
type CertificateSkkniResponse struct {
	UUID        uuid.UUID            `json:"uuid,omitempty"`
	Name        string               `json:"name,omitempty"`
	FileURL     string               `json:"file_url,omitempty"`
	Certificate *CertificateResponse `json:"certificate,omitempty"`
	CreatedAt   time.Time            `json:"created_at,omitempty"`
	UpdatedAt   time.Time            `json:"updated_at,omitempty"`
}

type SkkniThumbnailRequest struct {
	SkkniUUID uuid.UUID `validate:"required"`
}
