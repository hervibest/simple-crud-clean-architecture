package converter

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/model"
)

func CertifSkkniToResponse(certifSkkni *entity.Skkni) *model.CertificateSkkniResponse {

	var URL string

	if certifSkkni.File != nil {
		URL = certifSkkni.File.Path
	}

	var certificateReponsePointer *model.CertificateResponse
	if certifSkkni.Certificate != nil {
		certificateReponse := DTOCertificateToResponse(*certifSkkni.Certificate)
		certificateReponsePointer = &certificateReponse
	}

	return &model.CertificateSkkniResponse{
		UUID:        certifSkkni.UUID,
		Name:        certifSkkni.Name,
		FileURL:     URL,
		Certificate: certificateReponsePointer,
		CreatedAt:   certifSkkni.CreatedAt,
		UpdatedAt:   certifSkkni.UpdatedAt,
	}
}
