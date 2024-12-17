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
	Email    string `json:"email" validate:"required,email,min=5,max=100"`
	Password string `json:"password" validate:"required,min=6,max=100"`
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

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type GetUserRequest struct {
	Email string `validate:"required,max=100"`
}

type VerifyUserRequest struct {
	Token string `validate:"required"`
}

type LogoutUserRequest struct {
	Email        string
	AccessToken  string
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type AccessTokenRequest struct {
	Token string `validate:"required"`
}

type UpdateUserRequest struct {
	Email    string
	Password string `json:"password,omitempty" validate:"max=100"`
}
