package converter

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/model"
)

func DiscountToResponse(discount *entity.Discount) *model.DiscountResponse {

	courseResponse := Map(discount.Courses, DTOCourseToResponse)

	return &model.DiscountResponse{
		UUID:          discount.UUID,
		Name:          discount.Name,
		Type:          discount.Type,
		Value:         discount.Value,
		IsActive:      discount.IsActive,
		StartActiveAt: discount.StartActiveAt,
		ValidUntil:    discount.ValidUntil,
		Courses:       courseResponse,
		CreatedAt:     discount.CreatedAt,
		UpdatedAt:     discount.UpdatedAt,
	}
}

// func DiscountUUIDToResponse(discount *entity.Discount) *model.DiscountUUIDresponse {
// 	return &model.DiscountUUIDresponse{
// 		UUID: discount.UUID,
// 	}
// }
