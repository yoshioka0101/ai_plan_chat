package http

import (
	"github.com/gin-gonic/gin"
	"github.com/yoshioka0101/ai_plan_chat/internal/http/handler"
	"github.com/yoshioka0101/ai_plan_chat/internal/middleware"
)

// Server は統合ハンドラー
type Server struct {
	*handler.HealthHandler
	*handler.TaskHandler
}

// NewServer は統合ハンドラーを作成します
func NewServer(healthHandler *handler.HealthHandler, taskHandler *handler.TaskHandler) *Server {
	return &Server{
		HealthHandler: healthHandler,
		TaskHandler:   taskHandler,
	}
}

// SetupRoutes はルーターをセットアップします
func SetupRoutes(server *Server) *gin.Engine {
	logger := middleware.NewLogger()

	r := gin.New()
	r.Use(middleware.Logger(logger))
	r.Use(gin.Recovery())

	// Health check
	r.GET("/health", server.HealthHandler.GetHealth)

	// Task endpoints
	tasks := r.Group("/tasks")
	{
		tasks.GET("", server.TaskHandler.GetTaskList)
		tasks.POST("", server.TaskHandler.CreateTask)
		tasks.GET("/:id", server.TaskHandler.GetTask)
		tasks.PUT("/:id", server.TaskHandler.UpdateTask)
		tasks.PATCH("/:id", server.TaskHandler.EditTask)
		tasks.DELETE("/:id", server.TaskHandler.DeleteTask)
	}

	return r
}
