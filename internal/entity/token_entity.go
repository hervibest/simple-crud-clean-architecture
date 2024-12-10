package entity

import "time"

// AccessToken represents an access token entity
type AccessToken struct {
	UserID    string
	Token     string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	ExpiresAt time.Time
}

type RefreshToken struct {
	UserID    string
	Token     string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	ExpiresAt time.Time
}
