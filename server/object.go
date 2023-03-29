package server

import "github.com/google/uuid"

type Object struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	UserID string    `json:"-"`
}

func NewObject(name, userID string) Object {
	return Object{
		ID:     uuid.New(),
		Name:   name,
		UserID: userID,
	}
}
