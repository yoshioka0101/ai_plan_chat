package main

import (
	"fmt"
	"log"

	"github.com/yoshioka0101/ai_plan_chat/internal/utils"
)

func main() {
	fmt.Println("🔐 JWT Secret Generator")
	fmt.Println()

	// JWT_SECRETを生成
	secret, err := utils.GenerateJWTSecret()
	if err != nil {
		log.Fatalf("Failed to generate JWT secret: %v", err)
	}

	fmt.Printf("Generated JWT_SECRET: %s\n", secret)
	fmt.Println()
	fmt.Println("To use this secret, run:")
	fmt.Printf("export JWT_SECRET=\"%s\"\n", secret)
	fmt.Println()
	fmt.Println("Or add it to your .env file:")
	fmt.Printf("JWT_SECRET=%s\n", secret)
}
