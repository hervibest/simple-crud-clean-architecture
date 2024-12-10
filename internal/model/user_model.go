package model

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	UUID         uuid.UUID  `json:"uuid,omitempty"`
	Email        string     `json:"email,omitempty"`
	AccessToken  string     `json:"access_token,omitempty"`
	RefreshToken string     `json:"refresh_token,omitempty"`
	CreatedAt    time.Time  `json:"created_at,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at,omitempty"`
	VerifiedAt   *time.Time `json:"verified_at,omitempty"`
}

type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type VerifyEmailUserRequest struct {
	Email string `json:"email" validate:"required,email,max=100"`
	Token string `json:"token" validate:"required"`
}

type ResendEmailUserRequest struct {
	Email string `json:"email" validate:"required,email,max=100"`
}

type SendResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email,max=100"`
	Token string `json:"token" validate:"required"`
}

type ValidateResetTokenRequest struct {
	Email string `json:"email" validate:"required,max=100"`
	Token string `json:"token" validate:"required"`
}

type ResetPasswordUserRequest struct {
	Email    string `json:"email" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
	Token    string `json:"token" validate:"required"`
}
