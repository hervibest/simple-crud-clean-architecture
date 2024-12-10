package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateCourseRequest struct {
	Name          string `json:"name" validate:"required,max=255"`
	Slug          string
	Description   string   `json:"description"`
	Price         float64  `json:"price"`
	Categories    []string `json:"categories"`
	CategoryUUIDs []uuid.UUID
}

type GetCourseRequest struct {
	UUID uuid.UUID `json:"uuid,omitempty" validate:"required"`
}

type UpdateCourseRequest struct {
	UUID          uuid.UUID `json:"uuid,omitempty" validate:"required"`
	Name          string    `json:"name" validate:"required,max=255"`
	Slug          string
	Description   string   `json:"description"`
	Price         float64  `json:"price"`
	IsActive      bool     `json:"is_active"`
	Categories    []string `json:"categories"`
	CategoryUUIDs []uuid.UUID
}

type CourseResponse struct {
	UUID        uuid.UUID `json:"uuid,omitempty"`
	Name        string    `json:"email,omitempty"`
	Slug        string    `json:"slug,omitempty"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price,omitempty"`
	IsActive    bool      `json:"is_active,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
