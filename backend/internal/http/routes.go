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
	*handler.InterpretationHandler
	*handler.InterpretationItemHandler
}

// NewServer は統合ハンドラーを作成します
func NewServer(healthHandler *handler.HealthHandler, taskHandler *handler.TaskHandler, authHandler *handler.AuthHandler, interpretationHandler *handler.InterpretationHandler, interpretationItemHandler *handler.InterpretationItemHandler) *Server {
	return &Server{
		HealthHandler:              healthHandler,
		TaskHandler:                taskHandler,
		AuthHandler:                authHandler,
		InterpretationHandler:      interpretationHandler,
		InterpretationItemHandler: interpretationItemHandler,
	}
}

// SetupRoutes はルーターをセットアップします
func SetupRoutes(server *Server, authMiddleware *middleware.AuthMiddleware) *gin.Engine {
	logger := middleware.NewLogger()

	r := gin.New()
	r.Use(middleware.Logger(logger))
	r.Use(gin.Recovery())

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "https://app.hubplanner-ai.click"} // Add production frontend URL
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
		tasks.Use(authMiddleware.RequireAuth())
		{
			tasks.GET("", server.TaskHandler.GetTaskList)
			tasks.POST("", server.TaskHandler.CreateTask)
			tasks.GET("/:id", server.TaskHandler.GetTask)
			tasks.PUT("/:id", server.TaskHandler.UpdateTask)
			tasks.PATCH("/:id", server.TaskHandler.EditTask)
			tasks.DELETE("/:id", server.TaskHandler.DeleteTask)
		}

		// Interpretation endpoints
		interpretations := v1.Group("/interpretations")
		interpretations.Use(authMiddleware.RequireAuth())
		{
			interpretations.POST("", server.InterpretationHandler.CreateInterpretation)
			interpretations.GET("", server.InterpretationHandler.ListInterpretations)
			interpretations.GET("/:id", server.InterpretationHandler.GetInterpretation)
			interpretations.GET("/:id/items", server.InterpretationItemHandler.GetInterpretationItemsByInterpretationID)
			interpretations.POST("/:id/approve-items", server.InterpretationItemHandler.ApproveMultipleItems)
		}

		// Interpretation Item endpoints
		items := v1.Group("/interpretation-items")
		items.Use(authMiddleware.RequireAuth())
		{
			items.GET("/:id", server.InterpretationItemHandler.GetInterpretationItem)
			items.PATCH("/:id", server.InterpretationItemHandler.UpdateInterpretationItem)
			items.POST("/:id/approve", server.InterpretationItemHandler.ApproveInterpretationItem)
		}
	}

	return r
}
