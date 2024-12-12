package converter

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/model"
)

func EmployeeToResponse(employee *entity.Employee) *model.EmployeeResponse {
	return &model.EmployeeResponse{
		UUID:      employee.UUID,
		Email:     employee.Email,
		CreatedAt: employee.CreatedAt,
		UpdatedAt: employee.UpdatedAt,
	}
}

func EmployeeToTokenResponse(employee *entity.Employee) *model.EmployeeResponse {
	return &model.EmployeeResponse{
		UUID:        employee.UUID,
		Email:       employee.Email,
		CreatedAt:   employee.CreatedAt,
		UpdatedAt:   employee.UpdatedAt,
		AccessToken: employee.AccessToken,
	}
}
