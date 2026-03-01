package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sibukixxx/travelist/api/internal/apperror"
	"github.com/sibukixxx/travelist/api/internal/handler"
	"github.com/sibukixxx/travelist/api/internal/usecase"
)

// --- Test Doubles ---

type stubRegistrar struct {
	result *usecase.RegisterResult
	err    error
}

func (s *stubRegistrar) Register(_ context.Context, req usecase.RegisterRequest) (*usecase.RegisterResult, error) {
	return s.result, s.err
}

type stubVerifier struct {
	err error
}

func (s *stubVerifier) Verify(_ context.Context, token string) error {
	return s.err
}

// --- Tests ---

func TestUserHandlerRegister(t *testing.T) {
	t.Run("returns 201 on successful registration", func(t *testing.T) {
		h := handler.NewUserHandler(
			&stubRegistrar{result: &usecase.RegisterResult{UserID: "usr_123", Email: "test@example.com"}},
			&stubVerifier{},
		)

		body := `{"email":"test@example.com","password":"securepass"}`
		req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.Register(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
		}

		var resp usecase.RegisterResult
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp.UserID != "usr_123" {
			t.Errorf("UserID = %q, want %q", resp.UserID, "usr_123")
		}
		if resp.Email != "test@example.com" {
			t.Errorf("Email = %q, want %q", resp.Email, "test@example.com")
		}
	})

	t.Run("returns 400 for invalid JSON body", func(t *testing.T) {
		h := handler.NewUserHandler(&stubRegistrar{}, &stubVerifier{})

		req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader("not json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.Register(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}

		var resp map[string]any
		json.NewDecoder(w.Body).Decode(&resp)
		if resp["code"] != "BAD_REQUEST" {
			t.Errorf("code = %q, want %q", resp["code"], "BAD_REQUEST")
		}
	})

	t.Run("returns 409 for conflict error", func(t *testing.T) {
		h := handler.NewUserHandler(
			&stubRegistrar{err: apperror.NewConflict("email already registered")},
			&stubVerifier{},
		)

		body := `{"email":"dup@example.com","password":"securepass"}`
		req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.Register(w, req)

		if w.Code != http.StatusConflict {
			t.Errorf("status = %d, want %d", w.Code, http.StatusConflict)
		}

		var resp map[string]any
		json.NewDecoder(w.Body).Decode(&resp)
		if resp["code"] != "CONFLICT" {
			t.Errorf("code = %q, want %q", resp["code"], "CONFLICT")
		}
	})

	t.Run("returns 405 for non-POST method", func(t *testing.T) {
		h := handler.NewUserHandler(&stubRegistrar{}, &stubVerifier{})

		req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
		w := httptest.NewRecorder()

		h.Register(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
		}
	})
}

func TestUserHandlerVerifyEmail(t *testing.T) {
	t.Run("returns 200 on successful verification", func(t *testing.T) {
		h := handler.NewUserHandler(&stubRegistrar{}, &stubVerifier{})

		req := httptest.NewRequest(http.MethodGet, "/api/users/verify?token=validtoken", nil)
		w := httptest.NewRecorder()

		h.VerifyEmail(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
		}

		var resp map[string]string
		json.NewDecoder(w.Body).Decode(&resp)
		if resp["message"] != "email verified" {
			t.Errorf("message = %q, want %q", resp["message"], "email verified")
		}
	})

	t.Run("returns 400 for invalid token", func(t *testing.T) {
		h := handler.NewUserHandler(
			&stubRegistrar{},
			&stubVerifier{err: apperror.NewBadRequestWithCode(apperror.InvalidToken, "invalid verification token")},
		)

		req := httptest.NewRequest(http.MethodGet, "/api/users/verify?token=badtoken", nil)
		w := httptest.NewRecorder()

		h.VerifyEmail(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}

		var resp map[string]any
		json.NewDecoder(w.Body).Decode(&resp)
		if resp["code"] != "INVALID_TOKEN" {
			t.Errorf("code = %q, want %q", resp["code"], "INVALID_TOKEN")
		}
	})

	t.Run("returns 400 for missing token parameter", func(t *testing.T) {
		h := handler.NewUserHandler(
			&stubRegistrar{},
			&stubVerifier{err: apperror.NewBadRequestWithCode(apperror.InvalidToken, "verification token is required")},
		)

		req := httptest.NewRequest(http.MethodGet, "/api/users/verify", nil)
		w := httptest.NewRecorder()

		h.VerifyEmail(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}
