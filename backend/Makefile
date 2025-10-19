.PHONY: help lint fmt test dev build openapi-gen gen-api gen-db openapi redoc install dev-full test-db migrate

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
	@echo "  gen-db      データベースからGo型を生成"
	@echo "  redoc       ReDoc でOpenAPI仕様を確認"
	@echo "  bundle      OpenAPIファイルを結合"
	@echo "  dev-full    コード生成 + ドキュメント + 開発サーバー"
	@echo "  openapi-gen OpenAPIからコードを生成（従来版）"
	@echo "  openapi     ReDoc を起動（redoc のエイリアス）"
	@echo "  ci          GitHub Actions ワークフローの雛形をコピー"
	@echo "  test-db     データベース接続テスト"
	@echo "  migrate     マイグレーション実行"
	@echo "  help        このヘルプを表示"

# 必要なツールをインストール
install:
	@echo "開発に必要なツールをインストールしています..."
	@echo "Go ツール:"
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Node.js ツール:"
	npm install -g http-server
	npm install -g @redocly/cli
	@echo ""
	@echo "✅ インストール完了"
	@echo "使い方: make gen-api → make gen-db → make redoc → make dev"

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
	./bin/golangci-lint run

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
	@npx @redocly/cli@1.34.2 bundle internal/api/openapi/openapi.yaml -o internal/api/openapi/openapi-bundle.yaml

.PHONY: redocly-split
redocly-split: # redocly で OpenAPI を複数ファイルに分割する
	@npx @redocly/cli@1.34.2 split internal/api/openapi-bundle.yaml --outDir internal/api/schemas

.PHONY: gen-api
gen-api: redocly-bundle ## コード生成（oapi-codegen）
	@cd internal/api/openapi && go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=./oapi-codegen.yaml ./openapi-bundle.yaml
	@go mod tidy

.PHONY: gen-db
gen-db: ## データベースからGo型を生成（bobgen）
	@echo "データベースからGo型を生成しています"
	@echo "DB_DSN環境変数が設定されていない場合は、デフォルト値を使用します"
	@bobgen mysql -c bobgen.yaml
	@go mod tidy
	@echo "✅ データベース型生成完了: internal/gen/db"

# OpenAPIからコードを生成（従来版）
openapi-gen:
	@echo "OpenAPIからコードを生成しています"
	@cd internal/api && oapi-codegen --config=./oapi-codegen.yaml ./openapi.yaml

.PHONY: redoc
redoc: ## redoc を起動する
	npx @redocly/cli@1.34.2 preview-docs -p 3001 internal/api/openapi/openapi.yaml


# 開発用統合コマンド
dev-full:
	@echo "統合開発環境を起動中..."
	@$(MAKE) --no-print-directory gen-api
	@$(MAKE) --no-print-directory gen-db
	@echo "APIサーバー: http://localhost:8080"
	@echo "ReDoc: 別ターミナルで 'make redoc' を実行"
	@$(MAKE) --no-print-directory dev

# GitHub Actions ワークフローの雛形コピー
ci:
	@if [ -z "$(FILENAME)" ]; then \
		echo "FILENAME が指定されていません。例: make ci FILENAME=deploy"; \
		exit 1; \
	fi
	@mkdir -p ../.github/workflows
	@if [ -f ../.github/workflows/$(FILENAME).yml ]; then \
		echo "../.github/workflows/$(FILENAME).yml は既に存在します"; \
		exit 1; \
	fi
	cp ../.github/workflows/ci.yml.example ../.github/workflows/$(FILENAME).yml
	@echo ".github/workflows/$(FILENAME).yml を作成しました"

# データベース接続テスト
test-db:
	@echo "🔍 データベース接続テスト実行中..."
	@go run cmd/test-db/test-db.go

# マイグレーション実行
migrate:
	@echo "🚀 マイグレーション実行中..."
	@go run cmd/migrate/migrate.go