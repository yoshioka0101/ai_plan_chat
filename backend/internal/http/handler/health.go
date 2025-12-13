package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yoshioka0101/ai_plan_chat/gen/api"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) GetHealth(c *gin.Context) {
	response := api.HealthResponse{
		Status: "ok",
	}

	c.JSON(http.StatusOK, response)
}

// RegisterRoutes はヘルスチェック関連のルートを登録します
func (h *HealthHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/health", h.GetHealth)
}
