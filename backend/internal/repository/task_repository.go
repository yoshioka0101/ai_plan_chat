package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/mysql"
	"github.com/stephenafamo/bob/dialect/mysql/sm"
	"github.com/yoshioka0101/ai_plan_chat/gen/models"
	"github.com/yoshioka0101/ai_plan_chat/internal/interfaces"
)

type taskRepository struct {
	db     bob.Executor
	logger *slog.Logger
}

// NewTaskRepository は新しいTaskRepositoryを生成します
func NewTaskRepository(db *sql.DB, logger *slog.Logger) interfaces.TaskRepository {
	return &taskRepository{
		db:     bob.NewDB(db),
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
