package integration

import (
	"simple-crud-clean-architecture/internal/model"
	"testing"
)

func TestRegister(t *testing.T) {
	requestBody := model.RegisterUserRequest{
		Email:    "hervi.nur.r@365.ugm.ac.id",
		Password: "Hervi12345!!",
	}
}
