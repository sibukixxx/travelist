package apperror_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/sibukixxx/travelist/api/internal/apperror"
)

func TestAppErrorImplementsError(t *testing.T) {
	err := apperror.NewBadRequest("bad input")
	if err.Error() != "bad input" {
		t.Errorf("Error() = %q, want %q", err.Error(), "bad input")
	}
}

func TestAppErrorUnwrap(t *testing.T) {
	cause := errors.New("root cause")
	err := apperror.NewInternal(cause)

	if !errors.Is(err, cause) {
		t.Error("expected errors.Is to find wrapped cause")
	}
}

func TestAppErrorAs(t *testing.T) {
	err := apperror.NewNotFound("user not found")

	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatal("expected errors.As to match *AppError")
	}
	if appErr.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d, want %d", appErr.StatusCode, http.StatusNotFound)
	}
	if appErr.ErrCode != apperror.NotFound {
		t.Errorf("ErrCode = %q, want %q", appErr.ErrCode, apperror.NotFound)
	}
}

func TestConstructors(t *testing.T) {
	tests := []struct {
		name       string
		err        *apperror.AppError
		wantStatus int
		wantCode   apperror.Code
		wantMsg    string
	}{
		{
			name:       "NewBadRequest",
			err:        apperror.NewBadRequest("invalid"),
			wantStatus: http.StatusBadRequest,
			wantCode:   apperror.BadRequest,
			wantMsg:    "invalid",
		},
		{
			name:       "NewNotFound",
			err:        apperror.NewNotFound("missing"),
			wantStatus: http.StatusNotFound,
			wantCode:   apperror.NotFound,
			wantMsg:    "missing",
		},
		{
			name:       "NewConflict",
			err:        apperror.NewConflict("duplicate"),
			wantStatus: http.StatusConflict,
			wantCode:   apperror.Conflict,
			wantMsg:    "duplicate",
		},
		{
			name:       "NewValidation",
			err:        apperror.NewValidation("field error"),
			wantStatus: http.StatusBadRequest,
			wantCode:   apperror.ValidationError,
			wantMsg:    "field error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.StatusCode != tt.wantStatus {
				t.Errorf("StatusCode = %d, want %d", tt.err.StatusCode, tt.wantStatus)
			}
			if tt.err.ErrCode != tt.wantCode {
				t.Errorf("ErrCode = %q, want %q", tt.err.ErrCode, tt.wantCode)
			}
			if tt.err.Message != tt.wantMsg {
				t.Errorf("Message = %q, want %q", tt.err.Message, tt.wantMsg)
			}
		})
	}
}

func TestNewInternalHidesOriginalError(t *testing.T) {
	cause := errors.New("database connection failed")
	err := apperror.NewInternal(cause)

	if err.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %d, want %d", err.StatusCode, http.StatusInternalServerError)
	}
	if err.ErrCode != apperror.Internal {
		t.Errorf("ErrCode = %q, want %q", err.ErrCode, apperror.Internal)
	}
	if err.Message != "internal server error" {
		t.Errorf("Message = %q, want %q", err.Message, "internal server error")
	}
	if err.Err != cause {
		t.Error("expected wrapped error to be the cause")
	}
}

func TestNewBadRequestWithCode(t *testing.T) {
	err := apperror.NewBadRequestWithCode(apperror.InvalidToken, "token is invalid")

	if err.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode = %d, want %d", err.StatusCode, http.StatusBadRequest)
	}
	if err.ErrCode != apperror.InvalidToken {
		t.Errorf("ErrCode = %q, want %q", err.ErrCode, apperror.InvalidToken)
	}
	if err.Message != "token is invalid" {
		t.Errorf("Message = %q, want %q", err.Message, "token is invalid")
	}
}
