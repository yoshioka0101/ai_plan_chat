.PHONY: help lint fmt test dev build openapi-gen swagger

# デフォルトターゲット - ヘルプを表示
help:
	@echo "利用可能なコマンド:"
	@echo "  dev         開発環境の立ち上げ"
	@echo "  test        テストの実行"
	@echo "  lint        コードの静的解析"
	@echo "  fmt         コードフォーマット"
	@echo "  build       アプリケーションのビルド"
	@echo "  openapi-gen OpenAPIからコードを生成"
	@echo "  swagger     Swagger UI を起動"
	@echo "  help        このヘルプを表示"

# 開発環境の立ち上げ
dev:
	@echo "開発環境を起動しています"
	@export PATH=$$PATH:$$(go env GOPATH)/bin && air

# テストの実行
test:
	@echo "テストを実行しています"
	go test -v ./...

# コードの静的解析
lint:
	@echo "コードの lint をしますs"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint がインストールされていません。go vet を使用します。"; \
		go vet ./...; \
	fi

# コードフォーマット
fmt:
	@echo "コードをフォーマットしています"
	go fmt ./...

# アプリケーションのビルド
build:
	@echo "アプリケーションをビルドしています"
	go build -o bin/api cmd/api/main.go
	@echo "ビルドが完了しました: bin/api"

# OpenAPIからコードを生成
openapi-gen:
	@echo "OpenAPIからコードを生成しています"
	@if command -v oapi-codegen >/dev/null 2>&1; then \
		oapi-codegen --config=./oapi-codegen.yaml ./openapi.yaml; \
	else \
		echo "oapi-codegen がインストールされていません。インストールしてください:"; \
		echo "go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest"; \
	fi

# Swagger UI を起動
swagger:
	@echo "Swagger UI を起動しています"
	@docker run -p 3000:8080 -v $(PWD):/tmp -e SWAGGER_JSON=/tmp/openapi.yaml swaggerapi/swagger-ui