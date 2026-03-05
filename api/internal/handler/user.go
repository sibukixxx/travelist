package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sibukixxx/travelist/api/internal/apperror"
	"github.com/sibukixxx/travelist/api/internal/usecase"
)

// Registrar is the interface for user registration.
type Registrar interface {
	Register(ctx context.Context, req usecase.RegisterRequest) (*usecase.RegisterResult, error)
}

// Verifier is the interface for email verification.
type Verifier interface {
	Verify(ctx context.Context, token string) error
}

// UserHandler handles HTTP requests for user operations.
type UserHandler struct {
	registrar Registrar
	verifier  Verifier
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(registrar Registrar, verifier Verifier) *UserHandler {
	return &UserHandler{registrar: registrar, verifier: verifier}
}

// Register handles POST /api/users.
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req usecase.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleError(w, apperror.NewBadRequest("invalid request body"))
		return
	}

	result, err := h.registrar.Register(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// VerifyEmail handles GET /api/users/verify?token=...
func (h *UserHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if err := h.verifier.Verify(r.Context(), token); err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "email verified"})
}
