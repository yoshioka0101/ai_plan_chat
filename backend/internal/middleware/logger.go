package middleware

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// NewLogger はHTTPリクエストログ用のslogロガーを作成します。
// 必要最小限の情報のみを出力するように設定されています。
//
// 設定内容:
//   - Infoレベル以上のログのみ出力
//   - JSON形式で構造化されたログ
//   - stdout/stderrに出力
//
// Returns:
//
//	*slog.Logger: 設定済みのslogロガー
func NewLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}

func Logger(l *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// トレースIDを生成してコンテキストに設定
		traceID := uuid.New().String()
		c.Set("trace_id", traceID)

		start := time.Now()
		c.Next()

		// レスポンス時間を計算
		duration := time.Since(start)

		// ログレベルを決定（4xx, 5xxはwarn、その他はinfo）
		status := c.Writer.Status()
		var level slog.Level
		if status >= 400 {
			level = slog.LevelWarn
		} else {
			level = slog.LevelInfo
		}

		// 必要最小限の情報のみログ出力（トレースID含む）
		l.Log(c.Request.Context(), level, "HTTP",
			slog.String("trace_id", traceID),
			slog.Int("status", status),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("duration", duration.String()),
		)
	}
}
