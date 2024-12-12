package entity

import (
	"time"

	"github.com/google/uuid"
)

// AccessToken represents an access token entity
type AccessToken struct {
	UserUUID  uuid.UUID
	Token     string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	ExpiresAt time.Time
}

type RefreshToken struct {
	UserUUID  uuid.UUID
	Token     string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	ExpiresAt time.Time
}

type EmployeeAccessToken struct {
	EmployeeUUID uuid.UUID
	Token        string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	ExpiresAt    time.Time
}
