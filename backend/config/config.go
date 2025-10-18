package config

import (
	"fmt"
	"os"
)

// Config アプリケーションの設定を管理する構造体
type Config struct {
	// サーバー設定
	Port string

	// データベース設定
	Database DatabaseConfig
}

// DatabaseConfig データベース接続設定
type DatabaseConfig struct {
	Driver   string `json:"db_driver"`
	Host     string `json:"db_host"`
	Port     string `json:"db_port"`
	User     string `json:"db_user"`
	Password string `json:"-"`
	Name     string `json:"db_name"`
	DSN      string `json:"-"`
}

// Load 環境変数から設定を読み込む
func Load() *Config {
	// ポート設定（デフォルト8080）
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	config := &Config{
		Port: port,

		Database: DatabaseConfig{
			Driver:   getEnv("DB_DRIVER", "mysql"),
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", ""),
		},
	}

	// DSNを構築
	config.Database.DSN = buildDSN(config.Database)

	return config
}



// getEnv 環境変数を取得し、デフォルト値を設定
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}


// buildDSN データベース接続文字列を構築
func buildDSN(db DatabaseConfig) string {
	// 既にDB_DSNが設定されている場合はそれを使用
	if dsn := os.Getenv("DB_DSN"); dsn != "" {
		return dsn
	}

	// 個別の設定からDSNを構築
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
	)
}
