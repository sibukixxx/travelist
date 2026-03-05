package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sibukixxx/travelist/api/internal/usecase"
)

type UserHandler struct {
	registrar *usecase.UserRegistrar
}

type registerUserRequest struct {
	Email string `json:"email"`
}

func NewUserHandler(registrar *usecase.UserRegistrar) *UserHandler {
	return &UserHandler{registrar: registrar}
}

// Register handles POST /api/users/register.
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req registerUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	user, err := h.registrar.Register(r.Context(), req.Email)
	if err != nil {
		if usecase.IsEmailAlreadyExistsError(err) {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "email already registered"})
			return
		}
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, user)
}
