package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yoshioka0101/ai_plan_chat/gen/api"
)

// AppError はアプリケーション内部のエラー情報を持つ
type AppError struct {
	Code       string
	Message    string
	HTTPStatus int
}

// Error は error インターフェースの実装
func (e *AppError) Error() string {
	return e.Message
}

// 共通エラー定義
var (
	ErrInvalidRequest = &AppError{
		Code:       "invalid_request",
		Message:    "Invalid request parameters",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrUnauthorized = &AppError{
		Code:       "unauthorized",
		Message:    "Authentication required",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrConfigurationError = &AppError{
		Code:       "configuration_error",
		Message:    "Server configuration error",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrGeminiInitError = &AppError{
		Code:       "gemini_init_error",
		Message:    "Failed to initialize AI service",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrAIInterpretationError = &AppError{
		Code:       "ai_interpretation_error",
		Message:    "Failed to interpret input",
		HTTPStatus: http.StatusUnprocessableEntity,
	}

	ErrNotFound = &AppError{
		Code:       "not_found",
		Message:    "Resource not found",
		HTTPStatus: http.StatusNotFound,
	}

	ErrDatabaseError = &AppError{
		Code:       "database_error",
		Message:    "Database operation failed",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrInternalServer = &AppError{
		Code:       "internal_server_error",
		Message:    "An unexpected error occurred",
		HTTPStatus: http.StatusInternalServerError,
	}
)

// RespondWithError はエラーレスポンスを返す共通関数
func RespondWithError(c *gin.Context, err *AppError, details ...string) {
	message := err.Message
	if len(details) > 0 {
		message = details[0]
	}

	c.JSON(err.HTTPStatus, api.ErrorResponse{
		Code:    ptrString(err.Code),
		Message: ptrString(message),
	})
}

// RespondWithCustomError はカスタムエラーレスポンスを返す
func RespondWithCustomError(c *gin.Context, code string, message string, httpStatus int) {
	c.JSON(httpStatus, api.ErrorResponse{
		Code:    ptrString(code),
		Message: ptrString(message),
	})
}

// ptrString はstringのポインタを返すヘルパー関数
func ptrString(s string) *string {
	return &s
}
