package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yoshioka0101/ai_plan_chat/internal/gen/api"
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
