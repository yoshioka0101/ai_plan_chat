package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateJWTSecret は安全なJWT_SECRETを生成します
func GenerateJWTSecret() (string, error) {
	bytes := make([]byte, 32) // 32バイト = 256ビット
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// PrintJWTSecretGenerationHelp はJWT_SECRET生成のヘルプを表示します
func PrintJWTSecretGenerationHelp() {
	fmt.Println("=== JWT_SECRET Generation Help ===")
	fmt.Println("To generate a secure JWT_SECRET, run one of these commands:")
	fmt.Println()
	fmt.Println("1. Using OpenSSL:")
	fmt.Println("   openssl rand -base64 32")
	fmt.Println()
	fmt.Println("2. Using Go:")
	fmt.Println("   go run -c 'package main; import (\"crypto/rand\"; \"encoding/base64\"; \"fmt\"); func main() { bytes := make([]byte, 32); rand.Read(bytes); fmt.Println(base64.URLEncoding.EncodeToString(bytes)) }'")
	fmt.Println()
	fmt.Println("3. Set environment variable:")
	fmt.Println("   export JWT_SECRET=\"your-generated-secret-here\"")
	fmt.Println()

	// 例を生成
	if example, err := GenerateJWTSecret(); err == nil {
		fmt.Println("Generated JWT_SECRET example:", example)
	}
	fmt.Println("==================================")
}
