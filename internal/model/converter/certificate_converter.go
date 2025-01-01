package converter

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/model"
)

func CertificateToResponse(certificate *entity.Certificate) *model.CertificateResponse {

	finalPrice := certificate.Price

	// var discountResponsePointer *model.DiscountResponse

	// if len(certificate.Discounts) != 0 {
	// 	discountResponse := DTODiscountToResponse(certificate.Discounts[0])
	// 	discountResponsePointer = &discountResponse
	// 	if discountResponse.UUID != uuid.Nil && discountResponse.IsActive {
	// 		finalPrice = applyDiscount(certificate.Price, &discountResponse)
	// 	}
	// }

	var URL string

	if certificate.Thumbnail != nil {
		URL = certificate.Thumbnail.Path
	}

	var categoryReponsePointer *model.CertificateCatResponse
	if certificate.Category != nil {
		categoryReponse := DTOCertificateCatToResponse(*certificate.Category)
		categoryReponsePointer = &categoryReponse
	}

	return &model.CertificateResponse{
		UUID:        certificate.UUID,
		Name:        certificate.Name,
		Slug:        certificate.Slug,
		Description: certificate.Description,
		Price:       certificate.Price, // Harga asli
		FinalPrice:  finalPrice,        // Harga setelah diskon (atau harga asli jika tidak ada diskon)
		IsActive:    certificate.IsActive,
		// Discount:     discountResponsePointer,
		ThumbnailURL: URL,
		Category:     categoryReponsePointer,
		CreatedAt:    certificate.CreatedAt,
		UpdatedAt:    certificate.UpdatedAt,
	}
}

// func isDiscountValid(discount *model.DiscountResponse) bool {
// 	now := time.Now()
// 	return discount.ValidUntil.After(now) && discount.StartActiveAt.Before(now)
// }

// func applyDiscount(originalPrice float64, discount *model.DiscountResponse) float64 {
// 	if discount.Type == "PERCENT" {
// 		return originalPrice * (1 - discount.Value/100)
// 	} else {
// 		return originalPrice - discount.Value
// 	}
// }

func CertificateUUIDToResponse(certificate *entity.Certificate) *model.CertificateUUIDresponse {
	return &model.CertificateUUIDresponse{
		UUID: certificate.UUID,
	}
}

func DTOCertificateToResponse(certificateCat entity.Certificate) model.CertificateResponse {
	return model.CertificateResponse{
		UUID:        certificateCat.UUID,
		Name:        certificateCat.Name,
		Slug:        certificateCat.Slug,
		Description: certificateCat.Description,
		CreatedAt:   certificateCat.CreatedAt,
		UpdatedAt:   certificateCat.UpdatedAt,
	}
}
