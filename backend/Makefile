.PHONY: help lint fmt test dev build openapi-gen gen-api openapi redoc install dev-full

# デフォルトターゲット - ヘルプを表示
help:
	@echo "利用可能なコマンド:"
	@echo "  install 必要なツールをインストール"
	@echo "  dev         開発環境の立ち上げ"
	@echo "  test        テストの実行"
	@echo "  lint        コードの静的解析"
	@echo "  fmt         コードフォーマット"
	@echo "  build       アプリケーションのビルド"
	@echo "  gen-api     OpenAPIからコードを生成"
	@echo "  redoc       ReDoc でOpenAPI仕様を確認"
	@echo "  bundle      OpenAPIファイルを結合"
	@echo "  dev-full    コード生成 + ドキュメント + 開発サーバー"
	@echo "  openapi-gen OpenAPIからコードを生成（従来版）"
	@echo "  openapi     ReDoc を起動（redoc のエイリアス）"
	@echo "  help        このヘルプを表示"

# 必要なツールをインストール
install:
	@echo "開発に必要なツールをインストールしています..."
	@echo "Go ツール:"
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	go install github.com/air-verse/air@latest
	@echo "Node.js ツール:"
	npm install -g http-server
	npm install -g @redocly/cli
	@echo ""
	@echo "✅ インストール完了"
	@echo "使い方: make gen-api → make redoc → make dev"

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

.PHONY: redocly-bundle
redocly-bundle: # redocly で OpenAPI を単一ファイルにする
	@npx @redocly/cli@1.34.2 bundle openapi.yaml -o openapi-bundle.yaml

.PHONY: redocly-split
redocly-split: # redocly で OpenAPI を複数ファイルに分割する
	@npx @redocly/cli@1.34.2 split openapi-bundle.yaml --outDir schemas

.PHONY: gen-api
gen-api: redocly-bundle ## コード生成（oapi-codegen）
	@go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=./oapi-codegen.yaml ./openapi-bundle.yaml
	@go mod tidy

# OpenAPIからコードを生成（従来版）
openapi-gen:
	@echo "OpenAPIからコードを生成しています"
	@oapi-codegen --config=./oapi-codegen.yaml ./openapi.yaml

.PHONY: redoc
redoc: ## redoc を起動する
	npx @redocly/cli@1.34.2 preview-docs -p 3001 openapi.yaml


# 開発用統合コマンド
dev-full:
	@echo "統合開発環境を起動中..."
	@$(MAKE) --no-print-directory gen-api
	@echo "APIサーバー: http://localhost:8080"
	@echo "ReDoc: 別ターミナルで 'make redoc' を実行"
	@$(MAKE) --no-print-directory dev