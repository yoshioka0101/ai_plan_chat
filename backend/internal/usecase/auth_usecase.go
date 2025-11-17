package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yoshioka0101/ai_plan_chat/internal/repository"
	"github.com/yoshioka0101/ai_plan_chat/gen/models"
)

// AuthUsecase は認証ユースケースのインターフェースです
type AuthUsecase interface {
	SignUpWithGoogle(ctx context.Context, googleID, email, nickname, avatar string) (*models.User, error)
	SignInWithGoogle(ctx context.Context, googleID, email, nickname, avatar string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}

// authUsecase は認証ユースケースの実装です
type authUsecase struct {
	userRepo repository.UserRepository
}

// NewAuthUsecase は新しい認証ユースケースを作成します
func NewAuthUsecase(userRepo repository.UserRepository) AuthUsecase {
	return &authUsecase{
		userRepo: userRepo,
	}
}

// SignUpWithGoogle はGoogle認証で新規ユーザーを作成します
func (u *authUsecase) SignUpWithGoogle(ctx context.Context, googleID, email, nickname, avatar string) (*models.User, error) {
	// 既存ユーザーを検索
	existingUser, err := u.userRepo.GetUserByGoogleID(ctx, googleID)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		// 既存ユーザーが見つかった場合は、そのまま返す（新規作成ではないため）
		// 必要に応じて情報を更新
		now := time.Now()
		existingUser.Email = email
		existingUser.Nickname = nickname
		if avatar != "" {
			existingUser.Avatar.Set(avatar)
		}
		existingUser.UpdatedAt = now

		err = u.userRepo.UpdateUser(ctx, existingUser)
		if err != nil {
			return nil, err
		}
		return existingUser, nil
	}

	// 新規ユーザーを作成
	now := time.Now()
	newUser := &models.User{
		ID:        uuid.New().String(),
		GoogleID:  googleID,
		Email:     email,
		Nickname:  nickname,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if avatar != "" {
		newUser.Avatar.Set(avatar)
	}

	err = u.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

// SignInWithGoogle はGoogle認証で既存ユーザーを認証または更新します
// ユーザーが存在しない場合は新規作成します
func (u *authUsecase) SignInWithGoogle(ctx context.Context, googleID, email, nickname, avatar string) (*models.User, error) {
	// 既存ユーザーを検索
	existingUser, err := u.userRepo.GetUserByGoogleID(ctx, googleID)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	if existingUser != nil {
		// 既存ユーザーの情報を更新
		existingUser.Email = email
		existingUser.Nickname = nickname
		if avatar != "" {
			existingUser.Avatar.Set(avatar)
		}
		existingUser.UpdatedAt = now

		err = u.userRepo.UpdateUser(ctx, existingUser)
		if err != nil {
			return nil, err
		}
		return existingUser, nil
	}

	// ユーザーが見つからない場合は新規作成（相互依存を避けるため直接実装）
	newUser := &models.User{
		ID:        uuid.New().String(),
		GoogleID:  googleID,
		Email:     email,
		Nickname:  nickname,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if avatar != "" {
		newUser.Avatar.Set(avatar)
	}

	err = u.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

// GetUserByID はIDでユーザーを取得します
func (u *authUsecase) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return u.userRepo.GetUserByID(ctx, id)
}
