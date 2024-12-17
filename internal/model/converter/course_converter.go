package converter

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/model"

	"github.com/google/uuid"
)

func CourseToResponse(course *entity.Course) *model.CourseResponse {
	courseCatResponses := Map(course.Categories, DTOCourseCatToResponse)

	finalPrice := course.Price

	var discountResponsePointer *model.DiscountResponse

	if len(course.Discounts) != 0 {
		discountResponse := DTODiscountToResponse(course.Discounts[0])
		discountResponsePointer = &discountResponse
		if discountResponse.UUID != uuid.Nil && discountResponse.IsActive {
			finalPrice = applyDiscount(course.Price, &discountResponse)
		}
	}
	var URL string

	if course.Thumbnail != nil {
		URL = course.Thumbnail.Path
	}

	return &model.CourseResponse{
		UUID:         course.UUID,
		Name:         course.Name,
		Slug:         course.Slug,
		Description:  course.Description,
		Price:        course.Price, // Harga asli
		FinalPrice:   finalPrice,   // Harga setelah diskon (atau harga asli jika tidak ada diskon)
		IsActive:     course.IsActive,
		Discount:     discountResponsePointer,
		ThumbnailURL: URL,
		Categories:   courseCatResponses,
		CreatedAt:    course.CreatedAt,
		UpdatedAt:    course.UpdatedAt,
	}
}

// func isDiscountValid(discount *model.DiscountResponse) bool {
// 	now := time.Now()
// 	return discount.ValidUntil.After(now) && discount.StartActiveAt.Before(now)
// }

func applyDiscount(originalPrice float64, discount *model.DiscountResponse) float64 {
	if discount.Type == "PERCENT" {
		return originalPrice * (1 - discount.Value/100)
	} else {
		return originalPrice - discount.Value
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
