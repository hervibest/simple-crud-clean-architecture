package converter

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/model"
)

func CertifMaterialToResponse(certifMaterial *entity.Material) *model.CertificateMaterialResponse {

	return &model.CertificateMaterialResponse{
		UUID:      certifMaterial.UUID,
		Name:      certifMaterial.Name,
		Code:      certifMaterial.Code,
		Type:      certifMaterial.Type,
		CreatedAt: certifMaterial.CreatedAt,
		UpdatedAt: certifMaterial.UpdatedAt,
	}
}
