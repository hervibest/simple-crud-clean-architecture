package model

import (
	"time"

	"github.com/google/uuid"
)

type SearchCourseSecRequest struct {
	CourseID    int
	CourseUUID  uuid.UUID `json:"course_uuid" validate:"required"`
	Title       string
	Description string
	Page        int
	Size        int
}

type GetCourseSecRequest struct {
	UUID uuid.UUID `json:"uuid,omitempty" validate:"required"`
}

type CreateCourseSecRequest struct {
	CourseUUID  uuid.UUID `json:"course_uuid" validate:"required"`
	Title       string    `json:"title" validate:"required,max=255"`
	Sequence    int       `json:"sequence"`
	Description string    `json:"description"`
}

type UpdateCourseSecRequest struct {
	CourseUUID    uuid.UUID `json:"course_uuid" validate:"required"`
	CourseSecUUID uuid.UUID `validate:"required"`
	Title         string    `json:"title" validate:"required,max=255"`
	Sequence      int       `json:"sequence" validate:"required"`
	Description   string    `json:"description"`
}

type DeleteCourseSecRequest struct {
	CourseUUID    uuid.UUID `json:"course_uuid" validate:"required"`
	CourseSecUUID uuid.UUID `validate:"required"`
}
type CourseSectionResponse struct {
	UUID        uuid.UUID `json:"uuid,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Sequence    int       `json:"sequence,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
