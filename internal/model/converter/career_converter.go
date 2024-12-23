package converter

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/model"
)

func CareerToResponse(career *entity.Career) *model.CareerResponse {
	careerCatResponses := Map(career.Categories, DTOCareerCatToResponse)

	finalPrice := career.Price

	var URL string

	if career.Thumbnail != nil {
		URL = career.Thumbnail.Path
	}

	return &model.CareerResponse{
		UUID:         career.UUID,
		Title:        career.Title,
		Slug:         career.Slug,
		Description:  career.Description,
		Price:        career.Price, // Harga asli
		FinalPrice:   finalPrice,   // Harga setelah diskon (atau harga asli jika tidak ada diskon)
		IsActive:     career.IsActive,
		ThumbnailURL: URL,
		Categories:   careerCatResponses,
		CreatedAt:    career.CreatedAt,
		UpdatedAt:    career.UpdatedAt,
	}
}

// func isDiscountValid(discount *model.DiscountResponse) bool {
// 	now := time.Now()
// 	return discount.ValidUntil.After(now) && discount.StartActiveAt.Before(now)
// }

func CareerUUIDToResponse(career *entity.Career) *model.CareerUUIDresponse {
	return &model.CareerUUIDresponse{
		UUID: career.UUID,
	}
}

func DTOCareerToResponse(careerCat entity.Career) model.CareerResponse {
	return model.CareerResponse{
		UUID:        careerCat.UUID,
		Title:       careerCat.Title,
		Slug:        careerCat.Slug,
		Description: careerCat.Description,
		CreatedAt:   careerCat.CreatedAt,
		UpdatedAt:   careerCat.UpdatedAt,
	}
}
