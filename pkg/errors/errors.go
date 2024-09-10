// pkg/errors/errors.go

package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

type AppError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

func NewInternalServerError(message string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

func NewTooManyRequestsError(message string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusTooManyRequests,
	}
}

func RespondWithError(w http.ResponseWriter, err error) {
	var appErr *AppError
	var e *AppError
	if errors.As(err, &e) {
		appErr = e
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.StatusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": appErr.Message})
}
