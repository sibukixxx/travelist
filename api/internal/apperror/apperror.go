package apperror

import "net/http"

// Code represents a machine-readable error code.
type Code string

const (
	BadRequest      Code = "BAD_REQUEST"
	NotFound        Code = "NOT_FOUND"
	Conflict        Code = "CONFLICT"
	Internal        Code = "INTERNAL"
	InvalidToken    Code = "INVALID_TOKEN"
	ExpiredToken    Code = "EXPIRED_TOKEN"
	ValidationError Code = "VALIDATION_ERROR"
)

// AppError is a structured application error with HTTP status and error code.
type AppError struct {
	StatusCode int
	ErrCode    Code
	Message    string
	Err        error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewBadRequest(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		ErrCode:    BadRequest,
		Message:    message,
	}
}

func NewBadRequestWithCode(code Code, message string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		ErrCode:    code,
		Message:    message,
	}
}

func NewNotFound(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusNotFound,
		ErrCode:    NotFound,
		Message:    message,
	}
}

func NewConflict(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusConflict,
		ErrCode:    Conflict,
		Message:    message,
	}
}

func NewInternal(err error) *AppError {
	return &AppError{
		StatusCode: http.StatusInternalServerError,
		ErrCode:    Internal,
		Message:    "internal server error",
		Err:        err,
	}
}

func NewValidation(message string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		ErrCode:    ValidationError,
		Message:    message,
	}
}
