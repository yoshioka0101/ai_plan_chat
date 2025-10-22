package usecase

import (
	"context"
	"log/slog"

	"github.com/yoshioka0101/ai_plan_chat/gen/models"
	"github.com/yoshioka0101/ai_plan_chat/internal/interfaces"
)

type taskUsecase struct {
	repo   interfaces.TaskRepository
	logger *slog.Logger
}

// NewTaskUsecase は新しいTaskUsecaseを生成します
func NewTaskUsecase(repo interfaces.TaskRepository, logger *slog.Logger) interfaces.TaskUsecase {
	return &taskUsecase{
		repo:   repo,
		logger: logger,
	}
}

// GetTask はIDでタスクを取得します
func (u *taskUsecase) GetTask(ctx context.Context, id string) (*models.Task, error) {
	task, err := u.repo.GetTaskByID(ctx, id)
	if err != nil {
		u.logger.ErrorContext(ctx, "UseCase: Failed to get task",
			slog.String("task_id", id),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return task, nil
}
