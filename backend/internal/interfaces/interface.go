package interfaces

import (
	"context"
	"time"

	"github.com/yoshioka0101/ai_plan_chat/gen/models"
	"github.com/yoshioka0101/ai_plan_chat/internal/entity"
)

// TaskRepository はタスクのデータアクセスを提供します
type TaskRepository interface {
	GetTaskByID(ctx context.Context, id string) (*models.Task, error)
	GetAllTasks(ctx context.Context) (models.TaskSlice, error)
	CreateTask(ctx context.Context, task *models.Task) error
	UpdateTask(ctx context.Context, task *models.Task) error
	EditTask(ctx context.Context, id string, updates map[string]interface{}) (*models.Task, error)
	DeleteTask(ctx context.Context, id string) error
}

// TaskUsecase はタスクのビジネスロジックを提供します
type TaskUsecase interface {
	GetTask(ctx context.Context, id string) (*models.Task, error)
	GetTaskList(ctx context.Context) (models.TaskSlice, error)
	CreateTask(ctx context.Context, title string, description *string, dueAt *time.Time, status string) (*models.Task, error)
	UpdateTask(ctx context.Context, id string, title string, description *string, dueAt *time.Time, status string) (*models.Task, error)
	EditTask(ctx context.Context, id string, title *string, description *string, dueAt *time.Time, status *string) (*models.Task, error)
	DeleteTask(ctx context.Context, id string) error
}

// InterpretationRepository はAI解釈のデータアクセスを提供します
type InterpretationRepository interface {
	CreateInterpretation(ctx context.Context, interpretation *entity.AIInterpretation) error
	GetInterpretationByID(ctx context.Context, id string) (*entity.AIInterpretation, error)
	GetInterpretationsByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.AIInterpretation, error)
}
