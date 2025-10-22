package interfaces

import (
	"context"

	"github.com/yoshioka0101/ai_plan_chat/gen/models"
)

// TaskRepository はタスクのデータアクセスを提供します
type TaskRepository interface {
	GetTaskByID(ctx context.Context, id string) (*models.Task, error)
}

// TaskUsecase はタスクのビジネスロジックを提供します
type TaskUsecase interface {
	GetTask(ctx context.Context, id string) (*models.Task, error)
}
