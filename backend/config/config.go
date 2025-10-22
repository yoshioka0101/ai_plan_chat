package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
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
	// .envファイルを読み込む（エラーは無視 - 環境変数が直接設定されている場合もあるため）
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}

	// 必須環境変数のチェック（DB_DSNが設定されていない場合）
	if os.Getenv("DB_DSN") == "" {
		requiredEnvs := []string{"DB_USER", "DB_PASSWORD", "DB_NAME"}
		for _, env := range requiredEnvs {
			if os.Getenv(env) == "" {
				log.Fatalf("Required environment variable %s is not set (or set DB_DSN)", env)
			}
		}
	}

	// ポート設定（デフォルト8080）
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// データベースドライバー（デフォルトmysql）
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "mysql"
	}

	// データベースホスト（デフォルト127.0.0.1）
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	// データベースポート（デフォルト3306）
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}

	// DSNを構築
	var dsn string
	if envDSN := os.Getenv("DB_DSN"); envDSN != "" {
		dsn = envDSN
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			host,
			dbPort,
			os.Getenv("DB_NAME"),
		)
	}

	config := &Config{
		Port: port,

		Database: DatabaseConfig{
			Driver:   driver,
			Host:     host,
			Port:     dbPort,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			DSN:      dsn,
		},
	}

	return config
}
