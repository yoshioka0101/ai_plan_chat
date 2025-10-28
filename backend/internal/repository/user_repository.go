package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/mysql"
	"github.com/stephenafamo/bob/dialect/mysql/sm"
	"github.com/stephenafamo/bob/dialect/mysql/um"
	"github.com/yoshioka0101/ai_plan_chat/models"
)

// UserRepository はユーザーリポジトリのインターフェースです
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByGoogleID(ctx context.Context, googleID string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}

// userRepository はユーザーリポジトリの実装です
type userRepository struct {
	db     bob.Executor
	logger *slog.Logger
}

// NewUserRepository は新しいユーザーリポジトリを作成します
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db:     bob.NewDB(db),
		logger: slog.Default(),
	}
}

// CreateUser は新しいユーザーを作成します
func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	r.logger.InfoContext(ctx, "Repository: CreateUser started",
		slog.String("user_id", user.ID),
		slog.String("email", user.Email),
	)

	// UUIDを生成
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// 現在時刻を設定
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := models.Users.Insert(
		&models.UserSetter{
			ID:        omit.From(user.ID),
			GoogleID:  omit.From(user.GoogleID),
			Email:     omit.From(user.Email),
			Nickname:  omit.From(user.Nickname),
			Avatar:    omitnull.FromNull(user.Avatar),
			CreatedAt: omit.From(user.CreatedAt),
			UpdatedAt: omit.From(user.UpdatedAt),
		},
	).Exec(ctx, r.db)

	if err != nil {
		r.logger.ErrorContext(ctx, "Repository: Failed to create user",
			slog.String("user_id", user.ID),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to create user: %w", err)
	}

	r.logger.InfoContext(ctx, "Repository: CreateUser completed",
		slog.String("user_id", user.ID),
	)
	return nil
}

// GetUserByGoogleID はGoogleIDでユーザーを取得します
func (r *userRepository) GetUserByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	r.logger.InfoContext(ctx, "Repository: GetUserByGoogleID started",
		slog.String("google_id", googleID),
	)

	user, err := models.Users.Query(
		sm.Where(models.Users.Columns.GoogleID.EQ(mysql.Arg(googleID))),
	).One(ctx, r.db)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.WarnContext(ctx, "Repository: User not found",
				slog.String("google_id", googleID),
			)
			return nil, nil
		}
		r.logger.ErrorContext(ctx, "Repository: Failed to query user",
			slog.String("google_id", googleID),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	r.logger.InfoContext(ctx, "Repository: GetUserByGoogleID completed",
		slog.String("user_id", user.ID),
	)
	return user, nil
}

// GetUserByID はIDでユーザーを取得します
func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	r.logger.InfoContext(ctx, "Repository: GetUserByID started",
		slog.String("user_id", id.String()),
	)

	user, err := models.Users.Query(
		sm.Where(models.Users.Columns.ID.EQ(mysql.Arg(id.String()))),
	).One(ctx, r.db)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.WarnContext(ctx, "Repository: User not found",
				slog.String("user_id", id.String()),
			)
			return nil, nil
		}
		r.logger.ErrorContext(ctx, "Repository: Failed to query user",
			slog.String("user_id", id.String()),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	r.logger.InfoContext(ctx, "Repository: GetUserByID completed",
		slog.String("user_id", user.ID),
	)
	return user, nil
}

// UpdateUser はユーザーを更新します
func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	r.logger.InfoContext(ctx, "Repository: UpdateUser started",
		slog.String("user_id", user.ID),
	)

	// 更新時刻を設定
	user.UpdatedAt = time.Now()

	setter := &models.UserSetter{
		Email:     omit.From(user.Email),
		Nickname:  omit.From(user.Nickname),
		Avatar:    omitnull.FromNull(user.Avatar),
		UpdatedAt: omit.From(user.UpdatedAt),
	}

	_, err := models.Users.Update(
		setter.UpdateMod(),
		um.Where(models.Users.Columns.ID.EQ(mysql.Arg(user.ID))),
	).Exec(ctx, r.db)

	if err != nil {
		r.logger.ErrorContext(ctx, "Repository: Failed to update user",
			slog.String("user_id", user.ID),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to update user: %w", err)
	}

	r.logger.InfoContext(ctx, "Repository: UpdateUser completed",
		slog.String("user_id", user.ID),
	)
	return nil
}
