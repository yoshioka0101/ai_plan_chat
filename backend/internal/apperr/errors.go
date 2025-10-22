package apperr

import (
	"net/http"
)

// AppError はアプリケーション固有のエラー
type AppError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Error はエラーメッセージを返します
func (e *AppError) Error() string {
	return e.Message
}

// NewError は新しいアプリケーションエラーを作成します
func NewError(status int, message string) *AppError {
	return &AppError{
		Status:  status,
		Message: message,
	}
}

// Task関連のエラー
var (
	// 400 Bad Request
	ErrTaskValidationError = NewError(
		http.StatusBadRequest,
		"Validation error",
	)

	// 404 Not Found
	ErrTaskNotFound = NewError(
		http.StatusNotFound,
		"Task not found",
	)

	// 409 Conflict
	ErrTaskAlreadyExists = NewError(
		http.StatusConflict,
		"Task already exists",
	)

	// 500 Internal Server Error
	ErrTaskInternalError = NewError(
		http.StatusInternalServerError,
		"Internal server error",
	)

	// 500 Internal Server Error - Create failed
	ErrTaskCreateFailed = NewError(
		http.StatusInternalServerError,
		"Failed to create task",
	)

	// 500 Internal Server Error - Update failed
	ErrTaskUpdateFailed = NewError(
		http.StatusInternalServerError,
		"Failed to update task",
	)

	// 500 Internal Server Error - Delete failed
	ErrTaskDeleteFailed = NewError(
		http.StatusInternalServerError,
		"Failed to delete task",
	)
)
