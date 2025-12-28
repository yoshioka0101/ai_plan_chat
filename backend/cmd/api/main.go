package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yoshioka0101/ai_plan_chat/config"
	"github.com/yoshioka0101/ai_plan_chat/internal/infrastructure"
	"github.com/yoshioka0101/ai_plan_chat/internal/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	// 設定を読み込み
	cfg := config.Load()

	// OpenTelemetry初期化
	otelShutdown := func(context.Context) error { return nil }
	otelShutdown, err := telemetry.Init(context.Background(), cfg)
	if err != nil {
		log.Printf("Failed to initialize OpenTelemetry: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := otelShutdown(ctx); err != nil {
			log.Printf("Failed to shutdown OpenTelemetry: %v", err)
		}
	}()

	// データベース接続を初期化
	db, err := infrastructure.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("✅ Connected to database successfully")
	log.Printf("Starting server on port %s", cfg.Port)

	// サーバーを初期化
	r := InitializeServer(db, cfg)

	// HTTPサーバーを作成
	addr := fmt.Sprintf(":%s", cfg.Port)
	handler := otelhttp.NewHandler(r, "http.server")
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Graceful shutdownのためのチャンネル
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// サーバーを別のgoroutineで起動
	go func() {
		log.Printf("🚀 Server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// シグナルを待機
	<-quit
	log.Println("🛑 Shutting down server...")

	// Graceful shutdownのタイムアウトを設定
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// サーバーを停止
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("✅ Server exited gracefully")
	}

	// データベース接続を閉じる
	if err := db.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	} else {
		log.Println("✅ Database connection closed")
	}
}
