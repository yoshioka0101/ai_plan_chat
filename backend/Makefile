.PHONY: help lint fmt test dev build

# デフォルトターゲット - ヘルプを表示
help:
	@echo "利用可能なコマンド:"
	@echo "  dev     開発環境の立ち上げ"
	@echo "  test    テストの実行"
	@echo "  lint    コードの静的解析"
	@echo "  fmt     コードフォーマット"
	@echo "  build   アプリケーションのビルド"
	@echo "  help    このヘルプを表示"

# 開発環境の立ち上げ
dev:
	@echo "開発環境を起動しています"
	go run cmd/api/main.go

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