package presenter

import (
	"log"

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
	"github.com/yoshioka0101/ai_plan_chat/gen/api"
	"github.com/yoshioka0101/ai_plan_chat/gen/models"
)

// TaskPresenter はタスクのレスポンス整形を担当します
type TaskPresenter struct{}

func NewTaskPresenter() *TaskPresenter {
	return &TaskPresenter{}
}

// GetTask はBOBモデルをGetTask APIレスポンスに変換します
func (p *TaskPresenter) GetTask(task *models.Task) api.Task {
	// 文字列IDをUUIDにパース

	id, err := uuid.Parse(task.ID)
	if err != nil {
		// DB整合性が保たれていれば発生しないはず
		log.Printf("Warning: invalid UUID in database: %s, error: %v", task.ID, err)
		id = uuid.Nil
	}

	response := api.Task{
		Id:        types.UUID(id),
		Title:     task.Title,
		Status:    api.TaskStatus(task.Status),
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}

	// Null可能なフィールドの処理
	if val, ok := task.Description.Get(); ok {
		response.Description = &val
	}

	if val, ok := task.DueAt.Get(); ok {
		response.DueAt = &val
	}

	return response
}

// GetTaskList はBOBモデルスライスをGetTaskList APIレスポンスに変換します
func (p *TaskPresenter) GetTaskList(tasks models.TaskSlice) []api.Task {
	result := make([]api.Task, len(tasks))
	for i, task := range tasks {
		result[i] = p.GetTask(task)
	}
	return result
}

// CreateTask はBOBモデルをCreateTask APIレスポンスに変換します
func (p *TaskPresenter) CreateTask(task *models.Task) api.Task {
	return p.GetTask(task)
}

// UpdateTask はBOBモデルをUpdateTask APIレスポンスに変換します
func (p *TaskPresenter) UpdateTask(task *models.Task) api.Task {
	return p.GetTask(task)
}

// EditTask はBOBモデルをEditTask APIレスポンスに変換します
func (p *TaskPresenter) EditTask(task *models.Task) api.Task {
	return p.GetTask(task)
}
