package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/aarondl/opt/null"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/mysql"
	"github.com/stephenafamo/bob/dialect/mysql/dm"
	"github.com/stephenafamo/bob/dialect/mysql/sm"
	"github.com/stephenafamo/bob/dialect/mysql/um"
	"github.com/yoshioka0101/ai_plan_chat/gen/models"
	"github.com/yoshioka0101/ai_plan_chat/internal/interfaces"
)

type taskRepository struct {
	db     bob.Executor
	logger *slog.Logger
}

// NewTaskRepository は新しいTaskRepositoryを生成します
func NewTaskRepository(db *sql.DB, logger *slog.Logger) interfaces.TaskRepository {
	return NewTaskRepositoryWithExecutor(bob.NewDB(db), logger)
}

// NewTaskRepositoryWithExecutor は既存のexecutorを使ってTaskRepositoryを生成します
func NewTaskRepositoryWithExecutor(exec bob.Executor, logger *slog.Logger) interfaces.TaskRepository {
	return &taskRepository{
		db:     exec,
		logger: logger,
	}
}

// GetTaskByID はIDでタスクを取得します
func (r *taskRepository) GetTaskByID(ctx context.Context, id string) (*models.Task, error) {
	r.logger.InfoContext(ctx, "Repository: GetTaskByID started",
		slog.String("task_id", id),
	)

	task, err := models.Tasks.Query(
		sm.Where(models.Tasks.Columns.ID.EQ(mysql.Arg(id))),
	).One(ctx, r.db)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.WarnContext(ctx, "Repository: Task not found",
				slog.String("task_id", id),
			)
			return nil, fmt.Errorf("task not found: %s", id)
		}
		r.logger.ErrorContext(ctx, "Repository: Failed to query task",
			slog.String("task_id", id),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to find task: %w", err)
	}

	r.logger.InfoContext(ctx, "Repository: GetTaskByID completed",
		slog.String("task_id", id),
	)
	return task, nil
}

// GetAllTasks は全タスクを取得します
func (r *taskRepository) GetAllTasks(ctx context.Context) (models.TaskSlice, error) {
	r.logger.InfoContext(ctx, "Repository: GetAllTasks started")

	tasks, err := models.Tasks.Query(
		sm.OrderBy(mysql.Raw("created_at DESC")),
	).All(ctx, r.db)

	if err != nil {
		r.logger.ErrorContext(ctx, "Repository: Failed to query tasks",
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	r.logger.InfoContext(ctx, "Repository: GetAllTasks completed",
		slog.Int("count", len(tasks)),
	)
	return tasks, nil
}

// GetTasksByUserID はユーザーごとのタスク一覧を取得します
func (r *taskRepository) GetTasksByUserID(ctx context.Context, userID string) (models.TaskSlice, error) {
	r.logger.InfoContext(ctx, "Repository: GetTasksByUserID started",
		slog.String("user_id", userID),
	)

	tasks, err := models.Tasks.Query(
		sm.Where(models.Tasks.Columns.UserID.EQ(mysql.Arg(userID))),
		sm.OrderBy(mysql.Raw("created_at DESC")),
	).All(ctx, r.db)

	if err != nil {
		r.logger.ErrorContext(ctx, "Repository: Failed to query tasks by user",
			slog.String("user_id", userID),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	r.logger.InfoContext(ctx, "Repository: GetTasksByUserID completed",
		slog.String("user_id", userID),
		slog.Int("count", len(tasks)),
	)
	return tasks, nil
}

// CreateTask は新しいタスクを作成します
func (r *taskRepository) CreateTask(ctx context.Context, task *models.Task) error {
	r.logger.InfoContext(ctx, "Repository: CreateTask started",
		slog.String("task_id", task.ID),
		slog.String("title", task.Title),
	)

	// UUIDを生成
	if task.ID == "" {
		task.ID = uuid.New().String()
	}

	// 現在時刻を設定
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	// デフォルトステータスを設定
	if task.Status == "" {
		task.Status = "todo"
	}

	_, err := models.Tasks.Insert(
		&models.TaskSetter{
			ID:                 omit.From(task.ID),
			UserID:             omit.From(task.UserID),
			Title:              omit.From(task.Title),
			Description:        omitnull.FromNull(task.Description),
			DueAt:              omitnull.FromNull(task.DueAt),
			Status:             omit.From(task.Status),
			Source:             omit.From(task.Source),
			AiInterpretationID: omitnull.FromNull(task.AiInterpretationID),
			CreatedAt:          omit.From(task.CreatedAt),
			UpdatedAt:          omit.From(task.UpdatedAt),
		},
	).Exec(ctx, r.db)

	if err != nil {
		r.logger.ErrorContext(ctx, "Repository: Failed to create task",
			slog.String("task_id", task.ID),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to create task: %w", err)
	}

	r.logger.InfoContext(ctx, "Repository: CreateTask completed",
		slog.String("task_id", task.ID),
	)
	return nil
}

// UpdateTask はタスクを完全更新します
func (r *taskRepository) UpdateTask(ctx context.Context, task *models.Task) error {
	r.logger.InfoContext(ctx, "Repository: UpdateTask started",
		slog.String("task_id", task.ID),
	)

	// 更新時刻を設定
	task.UpdatedAt = time.Now()

	setter := &models.TaskSetter{
		Title:       omit.From(task.Title),
		Description: omitnull.FromNull(task.Description),
		DueAt:       omitnull.FromNull(task.DueAt),
		Status:      omit.From(task.Status),
		UpdatedAt:   omit.From(task.UpdatedAt),
	}

	_, err := models.Tasks.Update(
		setter.UpdateMod(),
		um.Where(models.Tasks.Columns.ID.EQ(mysql.Arg(task.ID))),
	).Exec(ctx, r.db)

	if err != nil {
		r.logger.ErrorContext(ctx, "Repository: Failed to update task",
			slog.String("task_id", task.ID),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to update task: %w", err)
	}

	r.logger.InfoContext(ctx, "Repository: UpdateTask completed",
		slog.String("task_id", task.ID),
	)
	return nil
}

// EditTask はタスクを部分更新します
func (r *taskRepository) EditTask(ctx context.Context, id string, updates map[string]interface{}) (*models.Task, error) {
	r.logger.InfoContext(ctx, "Repository: EditTask started",
		slog.String("task_id", id),
	)

	// 更新時刻を設定
	updates["updated_at"] = time.Now()

	// 更新用のSetterを作成
	setter := &models.TaskSetter{
		UpdatedAt: omit.From(updates["updated_at"].(time.Time)),
	}

	// 更新フィールドを設定
	if title, ok := updates["title"].(string); ok {
		setter.Title = omit.From(title)
	}
	if description, ok := updates["description"].(*string); ok {
		if description != nil {
			setter.Description = omitnull.FromNull(null.From(*description))
		}
	}
	if dueAt, ok := updates["due_at"].(*time.Time); ok {
		if dueAt != nil {
			setter.DueAt = omitnull.FromNull(null.From(*dueAt))
		}
	}
	if status, ok := updates["status"].(string); ok {
		setter.Status = omit.From(status)
	}

	_, err := models.Tasks.Update(
		setter.UpdateMod(),
		um.Where(models.Tasks.Columns.ID.EQ(mysql.Arg(id))),
	).Exec(ctx, r.db)

	if err != nil {
		r.logger.ErrorContext(ctx, "Repository: Failed to edit task",
			slog.String("task_id", id),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to edit task: %w", err)
	}

	// 更新されたタスクを取得
	task, err := r.GetTaskByID(ctx, id)
	if err != nil {
		r.logger.ErrorContext(ctx, "Repository: Failed to get updated task",
			slog.String("task_id", id),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to get updated task: %w", err)
	}

	r.logger.InfoContext(ctx, "Repository: EditTask completed",
		slog.String("task_id", id),
	)
	return task, nil
}

// DeleteTask はタスクを削除します
func (r *taskRepository) DeleteTask(ctx context.Context, id string) error {
	r.logger.InfoContext(ctx, "Repository: DeleteTask started",
		slog.String("task_id", id),
	)

	_, err := models.Tasks.Delete(
		dm.Where(models.Tasks.Columns.ID.EQ(mysql.Arg(id))),
	).Exec(ctx, r.db)

	if err != nil {
		r.logger.ErrorContext(ctx, "Repository: Failed to delete task",
			slog.String("task_id", id),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to delete task: %w", err)
	}

	r.logger.InfoContext(ctx, "Repository: DeleteTask completed",
		slog.String("task_id", id),
	)
	return nil
}
