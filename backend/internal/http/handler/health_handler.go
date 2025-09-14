package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) GetHealth(c *gin.Context) {
	// 外部スキーマ参照で生成された型を使用
	response := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	
	c.JSON(http.StatusOK, response)
}
