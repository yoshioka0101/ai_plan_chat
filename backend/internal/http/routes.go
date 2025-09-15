package http

import (
	"github.com/gin-gonic/gin"
	"github.com/yoshioka0101/ai_plane_chat/internal/http/handler"
	"github.com/yoshioka0101/ai_plane_chat/internal/middleware"
)

func SetupRoutes() *gin.Engine {
	logger := middleware.NewLogger()
	defer logger.Sync()

	r := gin.New()
	r.Use(middleware.Logger(logger))
	r.Use(gin.Recovery())

	// handlers
	healthHandler := handler.NewHealthHandler()

	// routes
	r.GET("/health", healthHandler.GetHealth)

	return r
}
