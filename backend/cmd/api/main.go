package main

import (
	"fmt"
	"log"

	"github.com/yoshioka0101/ai_plan_chat/config"
	"github.com/yoshioka0101/ai_plan_chat/internal/http"
)

func main() {
	// 設定を読み込み
	cfg := config.Load()
	
	// 設定情報をログ出力
	log.Printf("Starting server on port %s", cfg.Port)
	
	// ルーターをセットアップ
	r := http.SetupRoutes()
	
	// サーバーを起動
	addr := fmt.Sprintf(":%s", cfg.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
