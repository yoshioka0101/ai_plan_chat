package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yoshioka0101/ai_plan_chat/internal/http/handler"
	"github.com/yoshioka0101/ai_plan_chat/internal/middleware"
)

// Server は統合ハンドラー
type Server struct {
	*handler.HealthHandler
	*handler.TaskHandler
	*handler.AuthHandler
}

// NewServer は統合ハンドラーを作成します
func NewServer(healthHandler *handler.HealthHandler, taskHandler *handler.TaskHandler, authHandler *handler.AuthHandler) *Server {
	return &Server{
		HealthHandler: healthHandler,
		TaskHandler:   taskHandler,
		AuthHandler:   authHandler,
	}
}

// SetupRoutes はルーターをセットアップします
func SetupRoutes(server *Server) *gin.Engine {
	logger := middleware.NewLogger()

	r := gin.New()
	r.Use(middleware.Logger(logger))
	r.Use(gin.Recovery())

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"} // Vite dev server
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Health check
	r.GET("/health", server.HealthHandler.GetHealth)

	// Auth endpoints
	r.GET("/auth/google", server.AuthHandler.GoogleAuth)
	r.GET("/auth/google/callback", server.AuthHandler.GoogleCallback)
	r.POST("/auth/google/callback", server.AuthHandler.GoogleCallback)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Task endpoints
		tasks := v1.Group("/tasks")
		{
			tasks.GET("", server.TaskHandler.GetTaskList)
			tasks.POST("", server.TaskHandler.CreateTask)
			tasks.GET("/:id", server.TaskHandler.GetTask)
			tasks.PUT("/:id", server.TaskHandler.UpdateTask)
			tasks.PATCH("/:id", server.TaskHandler.EditTask)
			tasks.DELETE("/:id", server.TaskHandler.DeleteTask)
		}
	}

	return r
}
