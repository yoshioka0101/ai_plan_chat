package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yoshioka0101/ai_plan_chat/config"
)

func main() {
	// 設定を読み込み
	cfg := config.Load()

	// データベース接続文字列を設定から取得
	dsn := cfg.Database.DSN
	log.Printf("Using database DSN: %s", dsn)

	// データベースに接続
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// 接続確認
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("✅ Connected to database")

	// マイグレーションファイルのディレクトリ
	migrationsDir := "../migrations"

	// 絶対パスに変換
	absPath, err := filepath.Abs(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	log.Printf("Looking for migrations in: %s", absPath)

	// マイグレーションファイルを実行
	files, err := filepath.Glob(filepath.Join(absPath, "*.sql"))
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	if len(files) == 0 {
		log.Println("No migration files found")
		return
	}

	for _, file := range files {
		log.Printf("Executing migration: %s", filepath.Base(file))

		// ファイルを読み込み
		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", file, err)
		}

		// SQLを実行
		if _, err := db.Exec(string(content)); err != nil {
			log.Fatalf("Failed to execute migration %s: %v", file, err)
		}

		log.Printf("✅ Migration completed: %s", filepath.Base(file))
	}

	log.Println("✅ All migrations completed successfully")
}
