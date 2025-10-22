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
	// 404 Not Found
	ErrTaskNotFound = NewError(
		http.StatusNotFound,
		"Task not found",
	)

	// 500 Internal Server Error
	ErrTaskInternalError = NewError(
		http.StatusInternalServerError,
		"Internal server error",
	)
)
