package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config アプリケーションの設定を管理する構造体
type Config struct {
	// サーバー設定
	Port string

	// データベース設定
	Database DatabaseConfig

	// 認証設定
	Auth AuthConfig

	// AI設定
	AI AIConfig

	// Telemetry設定
	Telemetry TelemetryConfig
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

// AuthConfig 認証設定
type AuthConfig struct {
	JWTSecret          string `json:"-"`
	GoogleClientID     string `json:"-"`
	GoogleClientSecret string `json:"-"`
	GoogleRedirectURL  string `json:"-"`
}

// AIConfig AI設定
type AIConfig struct {
	GeminiAPIKey string `json:"-"`
	GeminiModel  string `json:"gemini_model"`
}

// TelemetryConfig OpenTelemetry設定
type TelemetryConfig struct {
	Enabled        bool
	ServiceName    string
	ServiceVersion string
	Environment    string
	OTLPEndpoint   string
	SampleRatio    float64
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

	// JWTシークレット（必須設定）
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		fmt.Println("❌ JWT_SECRET environment variable is not set!")
		fmt.Println()
		fmt.Println("=== JWT_SECRET Generation Help ===")
		fmt.Println("To generate a secure JWT_SECRET, run:")
		fmt.Println("   openssl rand -base64 32")
		fmt.Println()
		fmt.Println("Then set the environment variable:")
		fmt.Println("   export JWT_SECRET=\"your-generated-secret-here\"")
		fmt.Println("==================================")
		log.Fatal("JWT_SECRET environment variable is required. Please set a strong secret key.")
	}

	// JWT_SECRETの強度チェック（最低32文字）
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters long for security reasons.")
	}

	// GoogleOAuth設定（必須設定）
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if googleRedirectURL == "" {
		googleRedirectURL = "http://localhost:8080/auth/google/callback"
	}

	// GoogleOAuth設定の検証
	if googleClientID == "" || googleClientSecret == "" {
		log.Fatal("GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET environment variables are required for OAuth authentication.")
	}

	// Gemini API Key
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		log.Println("Warning: GEMINI_API_KEY is not set. AI features will not work.")
	}

	// Gemini Model
	// GEMINI_MODEL または GEMINI_MODEL_NAME から読み込み（後方互換性のため両方サポート）
	geminiModel := os.Getenv("GEMINI_MODEL")
	if geminiModel == "" {
		geminiModel = os.Getenv("GEMINI_MODEL_NAME")
	}
	if geminiModel == "" {
		// デフォルトモデル: Gemini 2.5 Flash-Lite (最も軽量で安価なモデル)
		geminiModel = "gemini-2.5-flash-lite"
		log.Printf("GEMINI_MODEL not set, using default: %s", geminiModel)
	}

	// OpenTelemetry settings
	otelEnabled := false
	if rawEnabled := os.Getenv("OTEL_ENABLED"); rawEnabled != "" {
		parsedEnabled, err := strconv.ParseBool(rawEnabled)
		if err != nil {
			log.Printf("Invalid OTEL_ENABLED value %q, defaulting to false", rawEnabled)
		} else {
			otelEnabled = parsedEnabled
		}
	}

	otelServiceName := os.Getenv("OTEL_SERVICE_NAME")
	if otelServiceName == "" {
		otelServiceName = "ai-plan-chat-api"
	}

	otelServiceVersion := os.Getenv("OTEL_SERVICE_VERSION")
	if otelServiceVersion == "" {
		otelServiceVersion = os.Getenv("APP_VERSION")
	}
	if otelServiceVersion == "" {
		otelServiceVersion = "unknown"
	}

	otelEnvironment := os.Getenv("OTEL_ENVIRONMENT")

	otelEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otelEndpoint == "" {
		otelEndpoint = "http://localhost:4318"
	}

	sampleRatio := 1.0
	if rawSampleRatio := os.Getenv("OTEL_SAMPLE_RATIO"); rawSampleRatio != "" {
		parsedRatio, err := strconv.ParseFloat(rawSampleRatio, 64)
		if err != nil {
			log.Printf("Invalid OTEL_SAMPLE_RATIO value %q, defaulting to 1.0", rawSampleRatio)
		} else {
			sampleRatio = parsedRatio
		}
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

		Auth: AuthConfig{
			JWTSecret:          jwtSecret,
			GoogleClientID:     googleClientID,
			GoogleClientSecret: googleClientSecret,
			GoogleRedirectURL:  googleRedirectURL,
		},

		AI: AIConfig{
			GeminiAPIKey: geminiAPIKey,
			GeminiModel:  geminiModel,
		},

		Telemetry: TelemetryConfig{
			Enabled:        otelEnabled,
			ServiceName:    otelServiceName,
			ServiceVersion: otelServiceVersion,
			Environment:    otelEnvironment,
			OTLPEndpoint:   otelEndpoint,
			SampleRatio:    sampleRatio,
		},
	}

	return config
}
