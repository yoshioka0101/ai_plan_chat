package main

import (
	"github.com/yoshioka0101/ai_plan_chat/internal/http"
)

func main() {
	r := http.SetupRoutes()
	r.Run(":8080")
}
