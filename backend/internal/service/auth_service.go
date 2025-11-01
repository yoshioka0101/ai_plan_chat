package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleoauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"

	"github.com/yoshioka0101/ai_plan_chat/models"
)

var (
	// ErrInvalidToken はトークンが無効な場合のエラーです
	ErrInvalidToken = errors.New("invalid token")
	// ErrInvalidState はstateパラメータが無効な場合のエラーです
	ErrInvalidState = errors.New("invalid state parameter")
	// ErrInvalidAlgorithm はJWTアルゴリズムが無効な場合のエラーです
	ErrInvalidAlgorithm = errors.New("invalid signing algorithm")
)

// OAuthState はOAuth2.0のstateパラメータを管理します
type OAuthState struct {
	State        string
	CodeVerifier string
	CreatedAt    time.Time
}

// AuthService は認証サービスのインターフェースです
type AuthService interface {
	// OAuth 2.0 + PKCE
	GenerateOAuthState() (*OAuthState, error)
	ValidateState(state string) (*OAuthState, error)
	GetAuthURL(state, codeChallenge string) string
	ExchangeGoogleCode(ctx context.Context, code, codeVerifier string) (*oauth2.Token, error)
	GetGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*googleoauth2.Userinfo, error)

	// JWT
	GenerateJWT(user *models.User) (string, error)
	ValidateJWT(tokenString string) (*jwt.Token, error)
}

// authService は認証サービスの実装です
type authService struct {
	jwtSecret    string
	googleConfig *oauth2.Config
	logger       *slog.Logger
	// CSRF対策: stateの一時保存 (本番環境ではRedis等を使用)
	stateStore map[string]*OAuthState
	stateMutex sync.RWMutex
}

// NewAuthService は新しい認証サービスを作成します
func NewAuthService(jwtSecret, googleClientID, googleClientSecret, redirectURL string) AuthService {
	config := &oauth2.Config{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	svc := &authService{
		jwtSecret:    jwtSecret,
		googleConfig: config,
		logger:       slog.Default(),
		stateStore:   make(map[string]*OAuthState),
	}

	// 古いstateを定期的にクリーンアップ (10分以上古いものを削除)
	go svc.cleanupExpiredStates()

	return svc
}

// cleanupExpiredStates は期限切れのstateを定期的に削除します
func (s *authService) cleanupExpiredStates() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.stateMutex.Lock()
		now := time.Now()
		for state, oauthState := range s.stateStore {
			if now.Sub(oauthState.CreatedAt) > 10*time.Minute {
				delete(s.stateStore, state)
				s.logger.Info("Expired OAuth state removed",
					slog.String("state", state),
				)
			}
		}
		s.stateMutex.Unlock()
	}
}

// GenerateOAuthState はOAuth2.0のstateとPKCE verifierを生成します (RFC 7636)
func (s *authService) GenerateOAuthState() (*OAuthState, error) {
	// 暗号学的に安全な乱数でstateを生成 (CSRF対策)
	stateBytes := make([]byte, 32)
	if _, err := rand.Read(stateBytes); err != nil {
		s.logger.Error("Failed to generate state", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to generate state: %w", err)
	}
	state := base64.URLEncoding.EncodeToString(stateBytes)

	// PKCE code_verifier生成 (RFC 7636 Section 4.1)
	verifierBytes := make([]byte, 32)
	if _, err := rand.Read(verifierBytes); err != nil {
		s.logger.Error("Failed to generate code verifier", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to generate code verifier: %w", err)
	}
	codeVerifier := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(verifierBytes)

	oauthState := &OAuthState{
		State:        state,
		CodeVerifier: codeVerifier,
		CreatedAt:    time.Now(),
	}

	// stateを一時保存 (本番環境ではRedis等の分散ストアを使用)
	s.stateMutex.Lock()
	s.stateStore[state] = oauthState
	s.stateMutex.Unlock()

	s.logger.Info("OAuth state generated", slog.String("state", state))
	return oauthState, nil
}

// ValidateState はstateパラメータを検証します (CSRF対策)
func (s *authService) ValidateState(state string) (*OAuthState, error) {
	s.stateMutex.RLock()
	oauthState, exists := s.stateStore[state]
	s.stateMutex.RUnlock()

	if !exists {
		s.logger.Warn("Invalid or expired state", slog.String("state", state))
		return nil, ErrInvalidState
	}

	// 10分以上経過したstateは無効
	if time.Since(oauthState.CreatedAt) > 10*time.Minute {
		s.stateMutex.Lock()
		delete(s.stateStore, state)
		s.stateMutex.Unlock()
		s.logger.Warn("Expired state", slog.String("state", state))
		return nil, ErrInvalidState
	}

	// 使用済みstateは削除 (replay attack対策)
	s.stateMutex.Lock()
	delete(s.stateStore, state)
	s.stateMutex.Unlock()

	s.logger.Info("State validated successfully", slog.String("state", state))
	return oauthState, nil
}

// GetAuthURL はOAuth認証URLを生成します (PKCE対応)
func (s *authService) GetAuthURL(state, codeChallenge string) string {
	// PKCE code_challenge生成 (RFC 7636 Section 4.2)
	// code_challenge = BASE64URL(SHA256(code_verifier))
	return s.googleConfig.AuthCodeURL(state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)
}

// GenerateCodeChallenge はcode_verifierからcode_challengeを生成します
func GenerateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hash[:])
}

// GenerateJWT はJWTトークンを生成します (RFC 7519準拠、有効期限1時間)
func (s *authService) GenerateJWT(user *models.User) (string, error) {
	now := time.Now()
	jti := uuid.New().String() // JWT ID (トークン失効用)

	claims := jwt.MapClaims{
		"jti":   jti,                           // JWT ID (RFC 7519 Section 4.1.7)
		"sub":   user.ID,                       // Subject (RFC 7519 Section 4.1.2)
		"email": user.Email,                    // カスタムクレーム
		"iat":   now.Unix(),                    // Issued At (RFC 7519 Section 4.1.6)
		"exp":   now.Add(1 * time.Hour).Unix(), // Expiration (RFC 7519 Section 4.1.4) - 1時間に短縮
		"nbf":   now.Unix(),                    // Not Before (RFC 7519 Section 4.1.5)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		s.logger.Error("Failed to generate JWT",
			slog.String("user_id", user.ID),
			slog.String("error", err.Error()),
		)
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	s.logger.Info("JWT generated",
		slog.String("user_id", user.ID),
		slog.String("jti", jti),
		slog.Time("exp", now.Add(1*time.Hour)),
	)
	return signedToken, nil
}

// ValidateJWT はJWTトークンを検証します (alg: none攻撃対策)
func (s *authService) ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// アルゴリズムのホワイトリスト検証 (alg: none攻撃対策)
		if token.Method != jwt.SigningMethodHS256 {
			s.logger.Warn("Invalid JWT algorithm detected",
				slog.String("alg", token.Method.Alg()),
			)
			return nil, ErrInvalidAlgorithm
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		s.logger.Warn("JWT validation failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if !token.Valid {
		s.logger.Warn("Invalid JWT token")
		return nil, ErrInvalidToken
	}

	// クレームの基本検証
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.Warn("Invalid JWT claims format")
		return nil, ErrInvalidToken
	}

	// 必須クレームの存在確認
	requiredClaims := []string{"sub", "exp", "iat", "jti"}
	for _, claim := range requiredClaims {
		if _, exists := claims[claim]; !exists {
			s.logger.Warn("Missing required JWT claim", slog.String("claim", claim))
			return nil, ErrInvalidToken
		}
	}

	s.logger.Info("JWT validated successfully",
		slog.String("user_id", claims["sub"].(string)),
		slog.String("jti", claims["jti"].(string)),
	)
	return token, nil
}

// ExchangeGoogleCode はGoogle認証コードをトークンに交換します (PKCE対応)
func (s *authService) ExchangeGoogleCode(ctx context.Context, code, codeVerifier string) (*oauth2.Token, error) {
	// PKCE code_verifierを使用してトークン交換 (RFC 7636 Section 4.5)
	token, err := s.googleConfig.Exchange(ctx, code,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier),
	)
	if err != nil {
		s.logger.Error("Failed to exchange Google code",
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	s.logger.Info("Google code exchanged successfully")
	return token, nil
}

// GetGoogleUserInfo はGoogleユーザー情報を取得します
func (s *authService) GetGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*googleoauth2.Userinfo, error) {
	client := s.googleConfig.Client(ctx, token)
	service, err := googleoauth2.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		s.logger.Error("Failed to create Google OAuth2 service",
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	userInfo, err := service.Userinfo.Get().Do()
	if err != nil {
		s.logger.Error("Failed to get Google user info",
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	s.logger.Info("Google user info retrieved",
		slog.String("email", userInfo.Email),
		slog.String("google_id", userInfo.Id),
	)
	return userInfo, nil
}
