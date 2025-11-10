package main

import (
	"database/sql"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/yoshioka0101/ai_plan_chat/config"
	"github.com/yoshioka0101/ai_plan_chat/internal/http"
	"github.com/yoshioka0101/ai_plan_chat/internal/http/handler"
	"github.com/yoshioka0101/ai_plan_chat/internal/http/presenter"
	"github.com/yoshioka0101/ai_plan_chat/internal/middleware"
	"github.com/yoshioka0101/ai_plan_chat/internal/repository"
	"github.com/yoshioka0101/ai_plan_chat/internal/service"
	"github.com/yoshioka0101/ai_plan_chat/internal/usecase"
)

// initializeHealthHandler はHealthHandlerを初期化します
func initializeHealthHandler() *handler.HealthHandler {
	return handler.NewHealthHandler()
}

// initializeTaskHandler はTaskHandlerとその依存関係を初期化します
func initializeTaskHandler(db *sql.DB, logger *slog.Logger) *handler.TaskHandler {
	// Repository → Usecase → Presenter → Handler
	taskRepo := repository.NewTaskRepository(db, logger)
	taskUsecase := usecase.NewTaskUsecase(taskRepo, logger)
	taskPresenter := presenter.NewTaskPresenter()
	return handler.NewTaskHandler(taskUsecase, taskPresenter)
}

// initializeAuthHandler はAuthHandlerとその依存関係を初期化します
func initializeAuthHandler(db *sql.DB, config *config.Config) *handler.AuthHandler {
	// Repository → Usecase → Service → Presenter → Handler
	userRepo := repository.NewUserRepository(db)
	authUsecase := usecase.NewAuthUsecase(userRepo)
	authService := service.NewAuthService(
		config.Auth.JWTSecret,
		config.Auth.GoogleClientID,
		config.Auth.GoogleClientSecret,
		config.Auth.GoogleRedirectURL,
	)
	authPresenter := presenter.NewAuthPresenter()
	return handler.NewAuthHandler(authUsecase, authService, authPresenter)
}

// initializeGeminiService はGeminiServiceを初期化します
func initializeGeminiService(config *config.Config) *service.GeminiService {
	if config.AI.GeminiAPIKey == "" {
		slog.Warn("Gemini API key is not set. AI features will not work.")
		return nil
	}

	geminiService, err := service.NewGeminiService(config.AI.GeminiAPIKey, config.AI.GeminiModel)
	if err != nil {
		slog.Error("Failed to initialize Gemini service", "error", err)
		return nil
	}

	slog.Info("Gemini service initialized successfully", "model", config.AI.GeminiModel)
	return geminiService
}

// initializeInterpretationHandler はInterpretationHandlerを初期化します
func initializeInterpretationHandler(db *sql.DB, logger *slog.Logger, geminiService *service.GeminiService) *handler.InterpretationHandler {
	interpretationRepo := repository.NewInterpretationRepository(db, logger)
	return handler.NewInterpretationHandler(geminiService, interpretationRepo)
}

// InitializeServer は全ての依存性注入を行い、Ginルーターを返します
func InitializeServer(db *sql.DB, config *config.Config) *gin.Engine {

	// Logger初期化
	logger := middleware.NewLogger()

	// サービスを初期化
	geminiService := initializeGeminiService(config)

	// 各ハンドラーを初期化
	healthHandler := initializeHealthHandler()
	taskHandler := initializeTaskHandler(db, logger)
	authHandler := initializeAuthHandler(db, config)
	interpretationHandler := initializeInterpretationHandler(db, logger, geminiService)

	// 統合ハンドラーを作成
	server := http.NewServer(healthHandler, taskHandler, authHandler, interpretationHandler)

	// ルーターをセットアップ（OpenAPI仕様に基づく）
	return http.SetupRoutes(server)
}
