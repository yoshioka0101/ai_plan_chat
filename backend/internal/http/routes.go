package http

import (
	"github.com/gin-gonic/gin"
	"github.com/yoshioka0101/ai_plan_chat/gen/api"
	"github.com/yoshioka0101/ai_plan_chat/internal/http/handler"
	"github.com/yoshioka0101/ai_plan_chat/internal/middleware"
)

// Server は api.ServerInterface を実装する統合ハンドラー
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

// AddTaskHandler はTaskHandlerを追加します（将来的な拡張用）
func (s *Server) AddTaskHandler(taskHandler *handler.TaskHandler) {
	s.TaskHandler = taskHandler
}

// SetupRoutes はルーターをセットアップします
func SetupRoutes(server *Server) *gin.Engine {
	logger := middleware.NewLogger()

	r := gin.New()
	r.Use(middleware.Logger(logger))
	r.Use(gin.Recovery())

	// OpenAPI仕様に基づいてルートを自動登録
	api.RegisterHandlers(r, server)

	return r
}
