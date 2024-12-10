package converter

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/model"
)

func CourseCatToResponse(courseCat *entity.CourseCategory) *model.CourseCatResponse {
	return &model.CourseCatResponse{
		UUID:        courseCat.UUID,
		Name:        courseCat.Name,
		Slug:        courseCat.Slug,
		Description: courseCat.Description,
		CreatedAt:   courseCat.CreatedAt,
		UpdatedAt:   courseCat.UpdatedAt,
	}
}
