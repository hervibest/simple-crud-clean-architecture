package converter

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/model"
)

func CareerCatToResponse(careerCat *entity.CareerCategory) *model.CareerCatResponse {
	return &model.CareerCatResponse{
		UUID:        careerCat.UUID,
		Name:        careerCat.Name,
		Slug:        careerCat.Slug,
		Description: careerCat.Description,
		CreatedAt:   careerCat.CreatedAt,
		UpdatedAt:   careerCat.UpdatedAt,
	}
}

func DTOCareerCatToResponse(careerCat entity.CareerCategory) model.CareerCatResponse {
	return model.CareerCatResponse{
		UUID:        careerCat.UUID,
		Name:        careerCat.Name,
		Slug:        careerCat.Slug,
		Description: careerCat.Description,
		CreatedAt:   careerCat.CreatedAt,
		UpdatedAt:   careerCat.UpdatedAt,
	}
}
