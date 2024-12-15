package model

import (
	"simple-crud-clean-architecture/internal/enum"
	"time"

	"github.com/google/uuid"
)

type CreateDiscountRequest struct {
	Name          string            `json:"name" validate:"required,max=255"`
	Type          enum.DiscountType `json:"type" validate:"required,max=100"`
	Value         float64           `json:"value"`
	StartActiveAt time.Time         `json:"start_active_at"`
	ValidUntil    time.Time         `json:"valid_until"`
	Courses       []string          `json:"course_uuid"`
	CourseUUIDs   []uuid.UUID
}

type GetDiscountRequest struct {
	UUID uuid.UUID `json:"uuid,omitempty" validate:"required"`
}

type UpdateDiscountRequest struct {
	UUID          uuid.UUID         `json:"uuid,omitempty" validate:"required"`
	Name          string            `json:"name" validate:"required,max=255"`
	Type          enum.DiscountType `json:"type" validate:"required,max=100"`
	Value         float64           `json:"value"`
	StartActiveAt time.Time         `json:"start_active_at"`
	ValidUntil    time.Time         `json:"valid_until"`
	Courses       []string          `json:"courses"`
	CourseUUIDs   []uuid.UUID
}

type ValidateDiscountRequest struct {
	UUID          uuid.UUID         `json:"uuid,omitempty"`
	Name          string            `json:"name"`
	Type          enum.DiscountType `json:"type"`
	Value         float64           `json:"value"`
	StartActiveAt time.Time         `json:"start_active_at"`
	ValidUntil    time.Time         `json:"valid_until"`
	Courses       []string          `json:"courses"`
	CourseUUIDs   []uuid.UUID
}

type DiscountResponse struct {
	UUID     uuid.UUID         `json:"uuid,omitempty"`
	Name     string            `json:"name"`
	Type     enum.DiscountType `json:"type"`
	Value    float64           `json:"value"`
	IsActive bool              `json:"is_active"`

	StartActiveAt time.Time        `json:"start_active_at"`
	ValidUntil    time.Time        `json:"valid_until"`
	Courses       []CourseResponse `json:"courses"`
	CreatedAt     time.Time        `json:"created_at,omitempty"`
	UpdatedAt     time.Time        `json:"updated_at,omitempty"`
}

type SearchDiscountRequest struct {
	Name string            `json:"name" validate:"max=255"`
	Type enum.DiscountType `json:"type" validate:"max=100"`
	Page int               `json:"page" validate:"min=1"`
	Size int               `json:"size" validate:"min=1,max=100"`
}

// type SearchPurchasedCourse struct {
// 	UserID      int    `validate:"required"`
// 	Name        string `json:"name" validate:"max=255"`
// 	Slug        string `json:"slug" validate:"max=255"`
// 	Description string `json:"description" validate:"max=255"`
// 	Page        int    `json:"page" validate:"min=1"`
// 	Size        int    `json:"size" validate:"min=1,max=100"`
// }

// type CourseUUIDresponse struct {
// 	UUID uuid.UUID `json:"uuid,omitempty"`
// }
