.PHONY: help lint fmt test dev build openapi-gen gen-api openapi redoc install dev-full gen-models migrate-diff migrate-apply db-sync

# デフォルトターゲット - ヘルプを表示
help:
	@echo "利用可能なコマンド:"
	@echo "  install        必要なツールをインストール"
	@echo "  dev            開発環境の立ち上げ"
	@echo "  test           テストの実行"
	@echo "  lint           コードの静的解析"
	@echo "  fmt            コードフォーマット"
	@echo "  build          アプリケーションのビルド"
	@echo ""
	@echo "コード生成:"
	@echo "  gen-api        OpenAPIからコードを生成"
	@echo "  gen-models     DBスキーマからGoの型を生成"
	@echo ""
	@echo "マイグレーション:"
	@echo "  migrate-diff   スキーマから差分マイグレーション生成 (make migrate-diff name=xxx)"
	@echo "  migrate-apply  マイグレーションをDBに適用"
	@echo "  db-sync        マイグレーション適用 + 型生成を一括実行"
	@echo ""
	@echo "ドキュメント:"
	@echo "  redoc          ReDoc でOpenAPI仕様を確認"
	@echo "  bundle         OpenAPIファイルを結合"
	@echo ""
	@echo "その他:"
	@echo "  dev-full       コード生成 + ドキュメント + 開発サーバー"
	@echo "  generate-jwt   JWT_SECRETを生成"
	@echo "  ci             GitHub Actions ワークフローの雛形をコピー"
	@echo "  help           このヘルプを表示"

# 必要なツールをインストール
install:
	@echo "開発に必要なツールをインストールしています..."
	@echo "Go ツール:"
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	go install github.com/air-verse/air@latest
	go install github.com/stephenafamo/bob/gen/bobgen-mysql@latest
	@echo "Node.js ツール:"
	npm install -g http-server
	npm install -g @redocly/cli
	@echo ""
	@echo "✅ インストール完了"

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
	@echo "コードの lint を実行しています..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run && echo "✅ リンターチェック完了 - エラーなし"; \
	else \
		echo "golangci-lint がインストールされていません。go vet を使用します。"; \
		go vet ./... && echo "✅ go vet チェック完了 - エラーなし"; \
	fi

# コードフォーマット
fmt:
	@echo "コードをフォーマットしています"
	go fmt ./...

# アプリケーションのビルド
build:
	@echo "アプリケーションをビルドしています"
	go build -o bin/api ./cmd/api
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

# =============================================================================
# Database & Code Generation
# =============================================================================

# DBスキーマからGoの型を生成
gen-models:
	@echo "データベーススキーマからGoの型を生成しています..."
	@export PATH=$$PATH:$$(go env GOPATH)/bin && bobgen-mysql
	@echo "不要なファイルをクリーンアップしています..."
	@find gen/models -name "bob_*" -type f -delete 2>/dev/null || true
	@go mod tidy
	@echo "✅ gen/models/ にテーブルモデルが生成されました"

# スキーマから差分マイグレーションを生成
migrate-diff:
	@if [ -z "$(name)" ]; then \
		echo "エラー: 変更内容を指定してください"; \
		echo "例: make migrate-diff name=create_users_table"; \
		exit 1; \
	fi
	@echo "マイグレーションを生成しています: $(name)"
	@cd ../migrations && atlas migrate diff $(name) --env dev
	@echo "✅ migrations/ にマイグレーションファイルが生成されました"

# マイグレーションをDBに適用
migrate-apply:
	@echo "マイグレーションをデータベースに適用しています..."
	@cd ../migrations && atlas migrate apply --env dev
	@echo "✅ マイグレーション適用完了"

# マイグレーション適用 + 型生成を一括実行
db-sync: migrate-apply gen-models
	@echo "✅ データベース同期完了"

# JWT_SECRETを生成
generate-jwt:
	@echo "JWT_SECRETを生成しています..."
	@go run cmd/generate-jwt-secret/main.go