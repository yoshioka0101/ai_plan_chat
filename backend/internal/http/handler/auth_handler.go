package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yoshioka0101/ai_plan_chat/internal/http/presenter"
	"github.com/yoshioka0101/ai_plan_chat/internal/service"
	"github.com/yoshioka0101/ai_plan_chat/internal/usecase"
)

// AuthHandler は認証ハンドラー
type AuthHandler struct {
	authUsecase usecase.AuthUsecase
	authService service.AuthService
	presenter   *presenter.AuthPresenter
	logger      *slog.Logger
}

// NewAuthHandler は新しい認証ハンドラーを作成
func NewAuthHandler(authUsecase usecase.AuthUsecase, authService service.AuthService, presenter *presenter.AuthPresenter) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		authService: authService,
		presenter:   presenter,
		logger:      slog.Default(),
	}
}

// GoogleCallbackRequest はGoogleOAuthコールバックのリクエスト
type GoogleCallbackRequest struct {
	Code  string `json:"code" binding:"required"`  // Authorization Code
	State string `json:"state" binding:"required"` // CSRF対策用state
}

// GoogleCallback はGoogleOAuthコールバックを処理します
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	ctx := c.Request.Context()
	var req GoogleCallbackRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WarnContext(ctx, "Invalid request body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Step 1: State検証
	oauthState, err := h.authService.ValidateState(req.State)
	if err != nil {
		if errors.Is(err, service.ErrInvalidState) {
			h.logger.WarnContext(ctx, "CSRF: Invalid state parameter",
				slog.String("state", req.State),
			)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired state parameter"})
			return
		}
		h.logger.ErrorContext(ctx, "State validation error", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication failed"})
		return
	}

	// Step 2: Authorization CodeをAccess Tokenに交換
	token, err := h.authService.ExchangeGoogleCode(ctx, req.Code, oauthState.CodeVerifier)
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to exchange code for token",
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	// Step 3: Googleユーザー情報を取得
	userInfo, err := h.authService.GetGoogleUserInfo(ctx, token)
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to get Google user info",
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	// Step 4: ユーザーを認証または作成
	user, err := h.authUsecase.SignInWithGoogle(
		ctx,
		userInfo.Id,
		userInfo.Email,
		userInfo.Name,
		userInfo.Picture,
	)
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to authenticate user",
			slog.String("google_id", userInfo.Id),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate user"})
		return
	}

	// Step 5: JWTトークンを生成
	jwtToken, err := h.authService.GenerateJWT(user)
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to generate JWT",
			slog.String("user_id", user.ID),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Step 6: 成功レスポンスを返す
	h.logger.InfoContext(ctx, "User authenticated successfully",
		slog.String("user_id", user.ID),
		slog.String("email", user.Email),
	)

	response := h.presenter.GoogleCallback(user, jwtToken)
	c.JSON(http.StatusOK, response)
}
