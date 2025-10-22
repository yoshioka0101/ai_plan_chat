package infrastructure

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yoshioka0101/ai_plan_chat/config"
)

// NewDB creates a new database connection
func NewDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.Database.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 接続プールの設定
	db.SetMaxOpenConns(25)                 // 最大オープン接続数
	db.SetMaxIdleConns(5)                  // 最大アイドル接続数
	db.SetConnMaxLifetime(5 * time.Minute) // 接続の最大生存時間

	// 接続確認
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
