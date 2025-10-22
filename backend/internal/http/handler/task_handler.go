package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yoshioka0101/ai_plan_chat/internal/apperr"
	"github.com/yoshioka0101/ai_plan_chat/internal/http/presenter"
	"github.com/yoshioka0101/ai_plan_chat/internal/interfaces"
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

// TODOあとで実装
func (h *TaskHandler) GetTaskList(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

func (h *TaskHandler) EditTask(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"})
}
