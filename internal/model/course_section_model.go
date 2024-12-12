package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateCourseSecRequest struct {
	CourseUUID  uuid.UUID `json:"course_uuid" validate:"required"`
	Title       string    `json:"title" validate:"required,max=255"`
	Sequence    int       `json:"sequence" validate:"required"`
	Description string    `json:"description"`
}

type CourseSectionResponse struct {
	UUID        uuid.UUID `json:"uuid,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Sequence    int       `json:"sequence,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
