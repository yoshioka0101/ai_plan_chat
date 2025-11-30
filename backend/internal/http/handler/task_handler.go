package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yoshioka0101/ai_plan_chat/gen/api"
	"github.com/yoshioka0101/ai_plan_chat/internal/apperr"
	"github.com/yoshioka0101/ai_plan_chat/internal/http/presenter"
	"github.com/yoshioka0101/ai_plan_chat/internal/interfaces"
	"github.com/yoshioka0101/ai_plan_chat/internal/validation"
)

// TaskHandler はタスク関連のHTTPハンドラー
type TaskHandler struct {
	usecase   interfaces.TaskUsecase
	presenter *presenter.TaskPresenter
}

func NewTaskHandler(usecase interfaces.TaskUsecase, presenter *presenter.TaskPresenter) *TaskHandler {
	return &TaskHandler{
		usecase:   usecase,
		presenter: presenter,
	}
}

// GetTask はタスクの単一取得 (GET /tasks/:id)
func (h *TaskHandler) GetTask(c *gin.Context) {
	ctx := c.Request.Context()
	taskID := c.Param("id")

	task, err := h.usecase.GetTask(ctx, taskID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			_ = c.Error(apperr.ErrTaskNotFound)
			return
		}
		_ = c.Error(apperr.ErrTaskInternalError)
		return
	}

	response := h.presenter.GetTask(task)
	c.JSON(http.StatusOK, response)
}

// GetTaskList はタスク一覧を取得します (GET /tasks)
func (h *TaskHandler) GetTaskList(c *gin.Context) {
	ctx := c.Request.Context()

	tasks, err := h.usecase.GetTaskList(ctx)
	if err != nil {
		_ = c.Error(apperr.ErrTaskInternalError)
		return
	}

	response := h.presenter.GetTaskList(tasks)
	c.JSON(http.StatusOK, response)
}

// CreateTask は新しいタスクを作成します (POST /tasks)
func (h *TaskHandler) CreateTask(c *gin.Context) {
	ctx := c.Request.Context()

	var req api.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperr.ErrTaskValidationError)
		return
	}

	// statusのデフォルト値設定
	status := "todo" // デフォルト値
	if req.Status != nil {
		status = string(*req.Status)
	}

	// バリデーション
	if err := validation.ValidateCreateTaskRequest(req.Title, req.Description, req.DueAt, status); err != nil {
		_ = c.Error(apperr.ErrTaskValidationError)
		return
	}

	task, err := h.usecase.CreateTask(ctx, req.Title, req.Description, req.DueAt, status)
	if err != nil {
		if strings.Contains(err.Error(), "validation") {
			_ = c.Error(apperr.ErrTaskValidationError)
			return
		}
		_ = c.Error(apperr.ErrTaskCreateFailed)
		return
	}

	response := h.presenter.CreateTask(task)
	c.JSON(http.StatusCreated, response)
}

// UpdateTask はタスクを完全更新します (PUT /tasks/:id)
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	ctx := c.Request.Context()
	taskID := c.Param("id")

	// IDのバリデーション
	if err := validation.ValidationTaskID(taskID); err != nil {
		_ = c.Error(apperr.ErrTaskValidationError)
		return
	}

	var req api.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperr.ErrTaskValidationError)
		return
	}

	// statusの検証（空文字列の場合はエラー）
	status := string(req.Status)
	if status == "" {
		_ = c.Error(apperr.ErrTaskValidationError)
		return
	}

	// バリデーション
	if err := validation.ValidateUpdateTaskRequest(req.Title, req.Description, req.DueAt, status); err != nil {
		_ = c.Error(apperr.ErrTaskValidationError)
		return
	}

	task, err := h.usecase.UpdateTask(ctx, taskID, req.Title, req.Description, req.DueAt, status)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			_ = c.Error(apperr.ErrTaskNotFound)
			return
		}
		if strings.Contains(err.Error(), "validation") {
			_ = c.Error(apperr.ErrTaskValidationError)
			return
		}
		_ = c.Error(apperr.ErrTaskUpdateFailed)
		return
	}

	response := h.presenter.UpdateTask(task)
	c.JSON(http.StatusOK, response)
}

// EditTask はタスクを部分更新します (PATCH /tasks/:id)
func (h *TaskHandler) EditTask(c *gin.Context) {
	ctx := c.Request.Context()
	taskID := c.Param("id")

	// IDのバリデーション
	if err := validation.ValidationTaskID(taskID); err != nil {
		_ = c.Error(apperr.ErrTaskValidationError)
		return
	}

	var req api.EditTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperr.ErrTaskValidationError)
		return
	}

	// バリデーション
	if err := validation.ValidateEditTaskRequest(req.Title, req.Description, req.DueAt, (*string)(req.Status)); err != nil {
		_ = c.Error(apperr.ErrTaskValidationError)
		return
	}

	// ステータスの型変換
	var status *string
	if req.Status != nil {
		s := string(*req.Status)
		status = &s
	}

	task, err := h.usecase.EditTask(ctx, taskID, req.Title, req.Description, req.DueAt, status)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			_ = c.Error(apperr.ErrTaskNotFound)
			return
		}
		if strings.Contains(err.Error(), "validation") {
			_ = c.Error(apperr.ErrTaskValidationError)
			return
		}
		_ = c.Error(apperr.ErrTaskUpdateFailed)
		return
	}

	response := h.presenter.EditTask(task)
	c.JSON(http.StatusOK, response)
}

// DeleteTask はタスクを削除します (DELETE /tasks/:id)
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	ctx := c.Request.Context()
	taskID := c.Param("id")

	// IDのバリデーション
	if err := validation.ValidationTaskID(taskID); err != nil {
		_ = c.Error(apperr.ErrTaskValidationError)
		return
	}

	err := h.usecase.DeleteTask(ctx, taskID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			_ = c.Error(apperr.ErrTaskNotFound)
			return
		}
		_ = c.Error(apperr.ErrTaskDeleteFailed)
		return
	}

	c.Status(http.StatusNoContent)
}
