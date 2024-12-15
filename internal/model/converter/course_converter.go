package converter

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/model"
)

func CourseToResponse(course *entity.Course) *model.CourseResponse {

	coureCatResponses := Map(course.Categories, DTOCourseCatToResponse)

	return &model.CourseResponse{
		UUID:        course.UUID,
		Name:        course.Name,
		Slug:        course.Slug,
		Description: course.Description,
		Price:       course.Price,
		IsActive:    course.IsActive,
		Categories:  coureCatResponses,
		CreatedAt:   course.CreatedAt,
		UpdatedAt:   course.UpdatedAt,
	}
}

func CourseUUIDToResponse(course *entity.Course) *model.CourseUUIDresponse {
	return &model.CourseUUIDresponse{
		UUID: course.UUID,
	}
}

func DTOCourseToResponse(courseCat entity.Course) model.CourseResponse {
	return model.CourseResponse{
		UUID:        courseCat.UUID,
		Name:        courseCat.Name,
		Slug:        courseCat.Slug,
		Description: courseCat.Description,
		CreatedAt:   courseCat.CreatedAt,
		UpdatedAt:   courseCat.UpdatedAt,
	}
}
