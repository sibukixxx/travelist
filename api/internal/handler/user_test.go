package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sibukixxx/travelist/api/internal/handler"
	"github.com/sibukixxx/travelist/api/internal/infra/clock"
	"github.com/sibukixxx/travelist/api/internal/infra/repo"
	"github.com/sibukixxx/travelist/api/internal/usecase"
)

func newUserHandler() *handler.UserHandler {
	userRepo := repo.NewInMemoryUserRepository()
	registrar := usecase.NewUserRegistrar(userRepo, clock.FixedClock{
		Time: time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC),
	})
	return handler.NewUserHandler(registrar)
}

func TestUserHandlerRegister(t *testing.T) {
	t.Run("returns 201 for valid email", func(t *testing.T) {
		h := newUserHandler()
		body := []byte(`{"email":"alice@example.com"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		h.Register(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("status = %d, want 201", rec.Code)
		}
		var got map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("invalid JSON response: %v", err)
		}
		if got["email"] != "alice@example.com" {
			t.Fatalf("email = %v, want alice@example.com", got["email"])
		}
	})

	t.Run("returns 409 for duplicated email", func(t *testing.T) {
		h := newUserHandler()
		body := []byte(`{"email":"alice@example.com"}`)

		firstReq := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader(body))
		firstRec := httptest.NewRecorder()
		h.Register(firstRec, firstReq)
		if firstRec.Code != http.StatusCreated {
			t.Fatalf("first status = %d, want 201", firstRec.Code)
		}

		secondReq := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader(body))
		secondRec := httptest.NewRecorder()
		h.Register(secondRec, secondReq)
		if secondRec.Code != http.StatusConflict {
			t.Fatalf("status = %d, want 409", secondRec.Code)
		}
	})

	t.Run("returns 400 for invalid email", func(t *testing.T) {
		h := newUserHandler()
		req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader([]byte(`{"email":"invalid"}`)))
		rec := httptest.NewRecorder()

		h.Register(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rec.Code)
		}
	})
}
