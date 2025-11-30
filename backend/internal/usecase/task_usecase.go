package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/aarondl/opt/null"
	"github.com/yoshioka0101/ai_plan_chat/gen/models"
	"github.com/yoshioka0101/ai_plan_chat/internal/interfaces"
	"github.com/yoshioka0101/ai_plan_chat/internal/validation"
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

// GetTaskList は全タスクを取得します
func (u *taskUsecase) GetTaskList(ctx context.Context) (models.TaskSlice, error) {
	u.logger.InfoContext(ctx, "UseCase: GetTaskList started")

	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		u.logger.WarnContext(ctx, "UseCase: Missing user_id in context for GetTaskList")
		return nil, fmt.Errorf("unauthorized")
	}

	tasks, err := u.repo.GetTasksByUserID(ctx, userID)
	if err != nil {
		u.logger.ErrorContext(ctx, "UseCase: Failed to get task list",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	u.logger.InfoContext(ctx, "UseCase: GetTaskList completed",
		slog.Int("count", len(tasks)),
	)
	return tasks, nil
}

// CreateTask は新しいタスクを作成します
func (u *taskUsecase) CreateTask(ctx context.Context, title string, description *string, dueAt *time.Time, status string) (*models.Task, error) {
	u.logger.InfoContext(ctx, "UseCase: CreateTask started",
		slog.String("title", title),
		slog.String("status", status),
	)

	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		u.logger.WarnContext(ctx, "UseCase: Missing user_id in context for CreateTask")
		return nil, fmt.Errorf("unauthorized")
	}

	// バリデーション
	if err := validation.ValidateCreateTaskRequest(title, description, dueAt, status); err != nil {
		u.logger.WarnContext(ctx, "UseCase: Validation failed",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	// デフォルトステータスを設定
	if status == "" {
		status = "todo"
	}

	// モデルを作成
	task := &models.Task{
		UserID: userID,
		Title:  title,
		Status: status,
		Source: "manual",
	}

	if description != nil {
		task.Description = null.From(*description)
	}
	if dueAt != nil {
		task.DueAt = null.From(*dueAt)
	}

	// リポジトリで作成
	if err := u.repo.CreateTask(ctx, task); err != nil {
		u.logger.ErrorContext(ctx, "UseCase: Failed to create task",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	u.logger.InfoContext(ctx, "UseCase: CreateTask completed",
		slog.String("task_id", task.ID),
	)
	return task, nil
}

// UpdateTask はタスクを完全更新します
func (u *taskUsecase) UpdateTask(ctx context.Context, id string, title string, description *string, dueAt *time.Time, status string) (*models.Task, error) {
	u.logger.InfoContext(ctx, "UseCase: UpdateTask started",
		slog.String("task_id", id),
		slog.String("title", title),
		slog.String("status", status),
	)

	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		u.logger.WarnContext(ctx, "UseCase: Missing user_id in context for UpdateTask")
		return nil, fmt.Errorf("unauthorized")
	}

	// バリデーション
	if err := validation.ValidateUpdateTaskRequest(title, description, dueAt, status); err != nil {
		u.logger.WarnContext(ctx, "UseCase: Validation failed",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	// 既存のタスクを取得
	existingTask, err := u.repo.GetTaskByID(ctx, id)
	if err != nil {
		u.logger.ErrorContext(ctx, "UseCase: Failed to get existing task",
			slog.String("task_id", id),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	if existingTask.UserID != userID {
		u.logger.WarnContext(ctx, "UseCase: Unauthorized task update attempt",
			slog.String("task_id", id),
			slog.String("user_id", userID),
		)
		return nil, fmt.Errorf("unauthorized")
	}

	// フィールドを更新
	existingTask.Title = title
	existingTask.Status = status

	if description != nil {
		existingTask.Description = null.From(*description)
	} else {
		existingTask.Description = null.Val[string]{}
	}

	if dueAt != nil {
		existingTask.DueAt = null.From(*dueAt)
	} else {
		existingTask.DueAt = null.Val[time.Time]{}
	}

	// リポジトリで更新
	if err := u.repo.UpdateTask(ctx, existingTask); err != nil {
		u.logger.ErrorContext(ctx, "UseCase: Failed to update task",
			slog.String("task_id", id),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	u.logger.InfoContext(ctx, "UseCase: UpdateTask completed",
		slog.String("task_id", id),
	)
	return existingTask, nil
}

// EditTask はタスクを部分更新します
func (u *taskUsecase) EditTask(ctx context.Context, id string, title *string, description *string, dueAt *time.Time, status *string) (*models.Task, error) {
	u.logger.InfoContext(ctx, "UseCase: EditTask started",
		slog.String("task_id", id),
	)

	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		u.logger.WarnContext(ctx, "UseCase: Missing user_id in context for EditTask")
		return nil, fmt.Errorf("unauthorized")
	}

	// バリデーション
	if err := validation.ValidateEditTaskRequest(title, description, dueAt, status); err != nil {
		u.logger.WarnContext(ctx, "UseCase: Validation failed",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	// 更新用のマップを作成
	updates := make(map[string]interface{})
	if title != nil {
		updates["title"] = *title
	}
	if description != nil {
		updates["description"] = description
	}
	if dueAt != nil {
		updates["due_at"] = dueAt
	}
	if status != nil {
		updates["status"] = *status
	}

	// リポジトリで部分更新
	task, err := u.repo.EditTask(ctx, id, updates)
	if err != nil {
		u.logger.ErrorContext(ctx, "UseCase: Failed to edit task",
			slog.String("task_id", id),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	u.logger.InfoContext(ctx, "UseCase: EditTask completed",
		slog.String("task_id", id),
	)
	return task, nil
}

// DeleteTask はタスクを削除します
func (u *taskUsecase) DeleteTask(ctx context.Context, id string) error {
	u.logger.InfoContext(ctx, "UseCase: DeleteTask started",
		slog.String("task_id", id),
	)

	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		u.logger.WarnContext(ctx, "UseCase: Missing user_id in context for DeleteTask")
		return fmt.Errorf("unauthorized")
	}

	// タスクの存在確認
	task, err := u.repo.GetTaskByID(ctx, id)
	if err != nil {
		u.logger.ErrorContext(ctx, "UseCase: Task not found",
			slog.String("task_id", id),
			slog.String("error", err.Error()),
		)
		return err
	}

	if task.UserID != userID {
		u.logger.WarnContext(ctx, "UseCase: Unauthorized delete attempt",
			slog.String("task_id", id),
			slog.String("user_id", userID),
		)
		return fmt.Errorf("unauthorized")
	}

	// リポジトリで削除
	if err := u.repo.DeleteTask(ctx, id); err != nil {
		u.logger.ErrorContext(ctx, "UseCase: Failed to delete task",
			slog.String("task_id", id),
			slog.String("error", err.Error()),
		)
		return err
	}

	u.logger.InfoContext(ctx, "UseCase: DeleteTask completed",
		slog.String("task_id", id),
	)
	return nil
}
