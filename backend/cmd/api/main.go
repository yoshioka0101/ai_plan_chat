package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yoshioka0101/ai_plane_chat/internal/middleware"
)

func main() {
	logger := middleware.NewLogger()
	defer logger.Sync()

	r := gin.New()

	// middleware として logger を設定
	r.Use(middleware.Logger(logger))
	// panic してもここで回復する｀
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.Run(":8080")
}
