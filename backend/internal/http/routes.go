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
		HealthHandler:             healthHandler,
		TaskHandler:               taskHandler,
		AuthHandler:               authHandler,
		InterpretationHandler:     interpretationHandler,
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

	// Explicitly handle OPTIONS for all routes as a fallback for CORS preflight
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(200)
	})

	// Register root level routes (health, auth)
	server.HealthHandler.RegisterRoutes(&r.RouterGroup)
	server.AuthHandler.RegisterRoutes(&r.RouterGroup)

	// API v1 routes with authentication
	v1 := r.Group("/api/v1")
	v1.Use(authMiddleware.RequireAuth())
	{
		server.TaskHandler.RegisterRoutes(v1)
		server.InterpretationHandler.RegisterRoutes(v1)
		server.InterpretationItemHandler.RegisterRoutes(v1)
	}

	return r
}
