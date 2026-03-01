package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sibukixxx/travelist/api/internal/apperror"
)

func TestWriteErrorJSON(t *testing.T) {
	t.Run("writes AppError as JSON with correct status", func(t *testing.T) {
		w := httptest.NewRecorder()
		appErr := apperror.NewBadRequest("invalid email")

		writeErrorJSON(w, appErr)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
		if ct := w.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want %q", ct, "application/json")
		}

		var resp ErrorResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp.Status != http.StatusBadRequest {
			t.Errorf("resp.Status = %d, want %d", resp.Status, http.StatusBadRequest)
		}
		if resp.Code != string(apperror.BadRequest) {
			t.Errorf("resp.Code = %q, want %q", resp.Code, apperror.BadRequest)
		}
		if resp.Message != "invalid email" {
			t.Errorf("resp.Message = %q, want %q", resp.Message, "invalid email")
		}
	})
}

func TestHandleError(t *testing.T) {
	t.Run("handles AppError correctly", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := apperror.NewConflict("email already exists")

		handleError(w, err)

		if w.Code != http.StatusConflict {
			t.Errorf("status = %d, want %d", w.Code, http.StatusConflict)
		}

		var resp ErrorResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp.Code != string(apperror.Conflict) {
			t.Errorf("resp.Code = %q, want %q", resp.Code, apperror.Conflict)
		}
	})

	t.Run("returns 500 with generic message for non-AppError", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := errors.New("database connection failed")

		handleError(w, err)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
		}

		var resp ErrorResponse
		if decErr := json.NewDecoder(w.Body).Decode(&resp); decErr != nil {
			t.Fatalf("failed to decode response: %v", decErr)
		}
		if resp.Code != string(apperror.Internal) {
			t.Errorf("resp.Code = %q, want %q", resp.Code, apperror.Internal)
		}
		if resp.Message != "internal server error" {
			t.Errorf("resp.Message = %q, want %q", resp.Message, "internal server error")
		}
	})
}
