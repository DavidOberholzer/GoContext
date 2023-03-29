package handler

import (
	"server"
	"server/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Store interface {
	Get(id uuid.UUID, userID string) (server.Object, error)
	Store(name, userID string) (server.Object, error)
}

type Handler struct {
	store  Store
	router chi.Router
}

func New(store Store) *Handler {
	h := &Handler{store: store}
	router := chi.NewRouter()

	router.Use(middleware.JSONResponse, middleware.TokenAuth)
	router.Post("/v1/objects", h.createObject)
	router.Get("/v1/objects/{id}", h.getObject)

	h.router = router

	return h
}
