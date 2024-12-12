package model

import (
	"time"

	"github.com/google/uuid"
)

type EmployeeResponse struct {
	UUID         uuid.UUID `json:"uuid,omitempty"`
	Email        string    `json:"email,omitempty"`
	AccessToken  string    `json:"access_token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

type RegisterEmployeeRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type LoginEmployeeRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type VerifyEmployeeRequest struct {
	Token string `validate:"required"`
}

type LogoutEmployeeRequest struct {
	Email       string
	AccessToken string
}

type GetEmployeeRequest struct {
	Email string `validate:"required,max=100"`
}
