package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"hashGenerationService/internal/model"
	"hashGenerationService/internal/service"
)

// HashService is the interface the handler depends on.
// Using an interface instead of *service.Service allows the handler to be
// tested with a mock without wiring up a real store.
type HashService interface {
	GenerateHash(input string) (*model.HashResponse, error)
	GetHash(hash string) (*model.HashResponse, error)
}

type Handler struct {
	service HashService
}

func NewHandler(svc HashService) *Handler {
	return &Handler{service: svc}
}

func (h *Handler) GenerateHash(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024)

	var req model.HashRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			respondError(w, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		respondError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if strings.TrimSpace(req.Input) == "" {
		respondError(w, http.StatusBadRequest, "input field is required")
		return
	}

	resp, err := h.service.GenerateHash(req.Input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			respondError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrMaxRetriesExceeded):
			respondError(w, http.StatusConflict, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondJSON(w, http.StatusCreated, resp)
}

func (h *Handler) GetHash(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("hash")
	if hash == "" {
		respondError(w, http.StatusBadRequest, "hash is required")
		return
	}

	resp, err := h.service.GetHash(hash)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrHashNotFound):
			respondError(w, http.StatusNotFound, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, model.ErrorResponse{Error: message})
}
