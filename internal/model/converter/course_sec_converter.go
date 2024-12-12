package converter

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/model"
)

func CourseSecToResponse(courseSec *entity.CourseSection) *model.CourseSectionResponse {

	return &model.CourseSectionResponse{
		UUID:        courseSec.UUID,
		Title:       courseSec.Title,
		Description: courseSec.Description,
		CreatedAt:   courseSec.CreatedAt,
		UpdatedAt:   courseSec.UpdatedAt,
	}
}
