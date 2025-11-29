package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yoshioka0101/ai_plan_chat/internal/service"
)

// AuthMiddleware はJWT認証ミドルウェアです
type AuthMiddleware struct {
	authService service.AuthService
	logger      *slog.Logger
}

// NewAuthMiddleware は新しい認証ミドルウェアを作成します
func NewAuthMiddleware(authService service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		logger:      slog.Default(),
	}
}

// RequireAuth はJWT認証を必須とするミドルウェアです
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Authorization ヘッダーの取得
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			m.logger.WarnContext(ctx, "Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Bearer トークンの抽出 (RFC 6750 Section 2.1)
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			m.logger.WarnContext(ctx, "Invalid Authorization header format",
				slog.String("header", authHeader),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Authorization header format. Expected: Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// JWTトークンの検証
		token, err := m.authService.ValidateJWT(tokenString)
		if err != nil {
			m.logger.WarnContext(ctx, "JWT validation failed",
				slog.String("error", err.Error()),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// クレームの取得
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			m.logger.ErrorContext(ctx, "Invalid JWT claims format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			c.Abort()
			return
		}

		// ユーザーIDの取得と検証
		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			m.logger.WarnContext(ctx, "Missing or invalid user ID in JWT claims")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			c.Abort()
			return
		}

		// コンテキストにユーザー情報を設定
		c.Set("user_id", userID)
		if email, ok := claims["email"].(string); ok {
			c.Set("email", email)
		}
		if jti, ok := claims["jti"].(string); ok {
			c.Set("jti", jti)
		}

		// Go標準のcontextにもユーザー情報を設定
		ctx = context.WithValue(ctx, "user_id", userID)
		if email, ok := claims["email"].(string); ok {
			ctx = context.WithValue(ctx, "email", email)
		}
		if jti, ok := claims["jti"].(string); ok {
			ctx = context.WithValue(ctx, "jti", jti)
		}
		c.Request = c.Request.WithContext(ctx)

		m.logger.InfoContext(ctx, "User authenticated",
			slog.String("user_id", userID),
		)

		c.Next()
	}
}

// OptionalAuth はJWT認証をオプションとするミドルウェアです
// 認証情報があれば検証しますが、なくてもリクエストを通します
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 認証情報なし: そのまま次へ
			c.Next()
			return
		}

		// Bearer トークンの抽出
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			// 不正な形式: ログのみ出力して次へ
			m.logger.WarnContext(ctx, "Invalid Authorization header format in optional auth",
				slog.String("header", authHeader),
			)
			c.Next()
			return
		}

		tokenString := parts[1]

		// JWTトークンの検証
		token, err := m.authService.ValidateJWT(tokenString)
		if err != nil {
			// 検証失敗: ログのみ出力して次へ
			m.logger.WarnContext(ctx, "JWT validation failed in optional auth",
				slog.String("error", err.Error()),
			)
			c.Next()
			return
		}

		// クレームの取得
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			m.logger.WarnContext(ctx, "Invalid JWT claims format in optional auth")
			c.Next()
			return
		}

		// ユーザーIDの取得
		if userID, ok := claims["sub"].(string); ok && userID != "" {
			c.Set("user_id", userID)
			if email, ok := claims["email"].(string); ok {
				c.Set("email", email)
			}
			if jti, ok := claims["jti"].(string); ok {
				c.Set("jti", jti)
			}

			// Go標準のcontextにもユーザー情報を設定
			ctx = context.WithValue(ctx, "user_id", userID)
			if email, ok := claims["email"].(string); ok {
				ctx = context.WithValue(ctx, "email", email)
			}
			if jti, ok := claims["jti"].(string); ok {
				ctx = context.WithValue(ctx, "jti", jti)
			}
			c.Request = c.Request.WithContext(ctx)

			m.logger.InfoContext(ctx, "User authenticated (optional)",
				slog.String("user_id", userID),
			)
		}

		c.Next()
	}
}
