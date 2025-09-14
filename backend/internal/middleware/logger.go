package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger はHTTPリクエストログ用のzapロガーを作成します。
// 必要最小限の情報のみを出力するように設定されています。
//
// 設定内容:
//   - Infoレベル以上のログのみ出力
//   - JSON形式で構造化されたログ
//   - timestamp, callerフィールドは出力しない（簡潔性のため）
//   - stdout/stderrに出力
//
// Returns:
//   *zap.Logger: 設定済みのzapロガー
func NewLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	// 不要なフィールドを無効化
	config.EncoderConfig.TimeKey = ""   // timestampを出力しない
	config.EncoderConfig.CallerKey = "" // callerを出力しない
	config.EncoderConfig.MessageKey = "msg"

	logger, _ := config.Build()
	return logger
}

func Logger(l *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		// レスポンス時間を計算
		duration := time.Since(start)

		// ログレベルを決定（4xx, 5xxはwarn、その他はinfo）
		status := c.Writer.Status()
		var level zapcore.Level
		if status >= 400 {
			level = zapcore.WarnLevel
		} else {
			level = zapcore.InfoLevel
		}

		// 必要最小限の情報のみログ出力
		l.Log(level, "HTTP",
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("duration", duration.String()), // duration.Stringで見やすく
		)
	}
}
