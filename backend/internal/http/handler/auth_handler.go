package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

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

// GoogleAuth はGoogleOAuth認証URLを生成します
func (h *AuthHandler) GoogleAuth(c *gin.Context) {
	ctx := c.Request.Context()

	// OAuth stateとPKCE verifierを生成
	oauthState, err := h.authService.GenerateOAuthState()
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to generate OAuth state", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize authentication"})
		return
	}

	// code_challengeを生成
	codeChallenge := service.GenerateCodeChallenge(oauthState.CodeVerifier)

	// Google認証URLを生成
	authURL := h.authService.GetAuthURL(oauthState.State, codeChallenge)

	h.logger.InfoContext(ctx, "Google auth URL generated", slog.String("state", oauthState.State))
	c.JSON(http.StatusOK, gin.H{"auth_url": authURL})
}

// googleCallbackRequest はGoogleOAuthコールバックのリクエスト（このパッケージ内でのみ使用）
type googleCallbackRequest struct {
	Code  string `form:"code" json:"code" binding:"required"`  // Authorization Code
	State string `form:"state" json:"state" binding:"required"` // CSRF対策用state
}

// GoogleCallback はGoogleOAuthコールバックを処理します
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	ctx := c.Request.Context()
	var req googleCallbackRequest

	// GETリクエストの場合はクエリパラメータから、POSTの場合はJSONから取得
	if c.Request.Method == http.MethodGet {
		if err := c.ShouldBindQuery(&req); err != nil {
			h.logger.WarnContext(ctx, "Invalid query parameters", slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
			return
		}
	} else {
		if err := c.ShouldBindJSON(&req); err != nil {
			h.logger.WarnContext(ctx, "Invalid request body", slog.String("error", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
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
		// Authorization codeの交換エラーは内部エラー（Google APIとの通信エラーなど）
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code for token"})
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
			slog.String("email", userInfo.Email),
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

	// Step 6: フロントエンドにリダイレクト（認証成功）
	h.logger.InfoContext(ctx, "User authenticated successfully",
		slog.String("user_id", user.ID),
		slog.String("email", user.Email),
	)

	// ユーザー情報をJSON形式に変換
	userData := map[string]interface{}{
		"id":       user.ID,
		"email":    user.Email,
		"nickname": user.Nickname,
	}
	if avatar, ok := user.Avatar.Get(); ok {
		userData["avatar"] = avatar
	}
	userJSONBytes, err := json.Marshal(userData)
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to marshal user data", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user data"})
		return
	}
	userJSON := url.QueryEscape(string(userJSONBytes))

	// フロントエンドのコールバックページにリダイレクト
	// tokenとuserをクエリパラメータとして渡す
	frontendURL := "http://localhost:5173/auth/callback"
	redirectURL := fmt.Sprintf("%s?token=%s&user=%s",
		frontendURL,
		url.QueryEscape(jwtToken),
		userJSON,
	)

	c.Redirect(http.StatusFound, redirectURL)
}

// RegisterRoutes は認証関連のルートを登録します
func (h *AuthHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/auth/google", h.GoogleAuth)
	group.GET("/auth/google/callback", h.GoogleCallback)
	group.POST("/auth/google/callback", h.GoogleCallback)
}
