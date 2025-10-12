package main

import (
	"log"

	"github.com/yoshioka0101/ai_plan_chat/internal/http"
)

func main() {
	r := http.SetupRoutes()
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
