package main

import (
	"fmt"
	"log"

	"github.com/yoshioka0101/ai_plan_chat/config"
	"github.com/yoshioka0101/ai_plan_chat/internal/infrastructure"
)

func main() {
	// 設定を読み込み
	cfg := config.Load()

	// データベース接続を初期化
	db, err := infrastructure.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	log.Println("✅ Connected to database successfully")
	log.Printf("Starting server on port %s", cfg.Port)

	// サーバーを初期化
	r := InitializeServer(db)

	// サーバーを起動
	addr := fmt.Sprintf(":%s", cfg.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
