package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/sibukixxx/travelist/api/internal/apperror"
)

// ErrorResponse is the standard JSON error response format.
type ErrorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeErrorJSON(w http.ResponseWriter, appErr *apperror.AppError) {
	writeJSON(w, appErr.StatusCode, ErrorResponse{
		Status:  appErr.StatusCode,
		Code:    string(appErr.ErrCode),
		Message: appErr.Message,
	})
}

func handleError(w http.ResponseWriter, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		if appErr.StatusCode == http.StatusInternalServerError {
			log.Printf("internal error: %v", appErr.Err)
		}
		writeErrorJSON(w, appErr)
		return
	}

	log.Printf("unexpected error: %v", err)
	writeErrorJSON(w, apperror.NewInternal(err))
}
