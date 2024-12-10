package converter

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/model"
)

func CourseToResponse(courseCat *entity.Course) *model.CourseResponse {
	return &model.CourseResponse{
		UUID:        courseCat.UUID,
		Name:        courseCat.Name,
		Slug:        courseCat.Slug,
		Description: courseCat.Description,
		Price:       courseCat.Price,
		IsActive:    courseCat.IsActive,
		CreatedAt:   courseCat.CreatedAt,
		UpdatedAt:   courseCat.UpdatedAt,
	}
}
