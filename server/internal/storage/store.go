package storage

import (
	"fmt"
	"server"

	"github.com/google/uuid"
)

type Store struct {
	db *db
}

func NewStore() *Store {
	return &Store{db: newDB()}
}

func (s *Store) Get(id uuid.UUID, userID string) (server.Object, error) {
	obj, err := s.db.Get(id.String())
	if err != nil {
		fmt.Println("Log Get Error here if unexpected")
		return server.Object{}, err
	}

	// Objects not owned by the user will report 404.
	if obj.UserID != userID {
		return server.Object{}, ErrNotFound
	}

	return obj, nil
}

func (s *Store) Store(name, userID string) (server.Object, error) {
	obj := server.NewObject(name, userID)

	if err := s.db.Set(obj.ID.String(), obj); err != nil {
		fmt.Println("Log Set Error here if unexpected")
		return server.Object{}, err
	}

	return obj, nil
}
