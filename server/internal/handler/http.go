package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"server/internal/middleware"
	"server/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s request on path - %s\n", r.Method, r.URL.Path)
	h.router.ServeHTTP(w, r)
}

// getObject retrieves an object from storage.
func (h *Handler) getObject(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromCtx(r.Context())
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	pathParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(pathParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	obj, err := h.store.Get(id, userID)
	switch {
	case errors.Is(err, storage.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(obj); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// createObject creates an object in storage.
func (h *Handler) createObject(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromCtx(r.Context())
	if userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var createRequest apiCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := createRequest.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	obj, err := h.store.Store(createRequest.Name, userID)
	switch {
	case errors.Is(err, storage.ErrAlreadyExists):
		w.WriteHeader(http.StatusConflict)
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(obj); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
