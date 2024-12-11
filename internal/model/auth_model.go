package model

import "github.com/google/uuid"

type Auth struct {
	// Login user id
	UUID  uuid.UUID
	Id    int
	Email string
	Token string
}
