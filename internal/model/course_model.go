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

type GetPurchasedCourseRequest struct {
	UserID int `validate:"required"`
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
	UUID         uuid.UUID           `json:"uuid,omitempty"`
	Name         string              `json:"name,omitempty"`
	Slug         string              `json:"slug,omitempty"`
	Description  string              `json:"description,omitempty"`
	Price        float64             `json:"original_price,omitempty"`
	FinalPrice   float64             `json:"final_price,omitempty"`
	IsActive     bool                `json:"is_active,omitempty"`
	Categories   []CourseCatResponse `json:"categories,omitempty"`
	ThumbnailURL string              `json:"thumbnail_url,omitempty"`
	Discount     *DiscountResponse   `json:"discount,omitempty"`
	CreatedAt    time.Time           `json:"created_at,omitempty"`
	UpdatedAt    time.Time           `json:"updated_at,omitempty"`
}

type SearchCourseRequest struct {
	Name        string `json:"name" validate:"max=255"`
	Slug        string `json:"slug" validate:"max=255"`
	Description string `json:"description" validate:"max=255"`
	Page        int    `json:"page" validate:"min=1"`
	Size        int    `json:"size" validate:"min=1,max=100"`
}

type SearchPurchasedCourse struct {
	UserID      int    `validate:"required"`
	Name        string `json:"name" validate:"max=255"`
	Slug        string `json:"slug" validate:"max=255"`
	Description string `json:"description" validate:"max=255"`
	Page        int    `json:"page" validate:"min=1"`
	Size        int    `json:"size" validate:"min=1,max=100"`
}

type CourseUUIDresponse struct {
	UUID uuid.UUID `json:"uuid,omitempty"`
}

type UploadThumbnailRequest struct {
	CourseUUID uuid.UUID `validate:"required"`
}
