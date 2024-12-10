package model

import "github.com/google/uuid"

type Auth struct {
	// Login user id
	UUID  uuid.UUID
	Email string
	Token string
}
