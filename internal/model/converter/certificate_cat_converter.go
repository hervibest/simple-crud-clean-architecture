package converter

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/model"
)

func CertificateCatToResponse(certificateCat *entity.CertificateCategory) *model.CertificateCatResponse {
	return &model.CertificateCatResponse{
		UUID:        certificateCat.UUID,
		Name:        certificateCat.Name,
		Slug:        certificateCat.Slug,
		Description: certificateCat.Description,
		CreatedAt:   certificateCat.CreatedAt,
		UpdatedAt:   certificateCat.UpdatedAt,
	}
}

func DTOCertificateCatToResponse(certificateCat entity.CertificateCategory) model.CertificateCatResponse {
	return model.CertificateCatResponse{
		UUID:        certificateCat.UUID,
		Name:        certificateCat.Name,
		Slug:        certificateCat.Slug,
		Description: certificateCat.Description,
		CreatedAt:   certificateCat.CreatedAt,
		UpdatedAt:   certificateCat.UpdatedAt,
	}
}
