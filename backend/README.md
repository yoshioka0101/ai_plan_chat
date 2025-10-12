# AI Plan Chat の Backend開発
Go言語で構築されたREST APIサーバーです。

## 技術スタック

- フレームワーク: [Gin](https://github.com/gin-gonic/gin)
- ログ: [zap](https://github.com/uber-go/zap) (構造化ログ)
- OpenAPI: [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) (コード生成)
- ホットリロード: [Air](https://github.com/air-verse/air)

## 必要なツール

- Go 1.25

## セットアップ

```bash
# 必要なツールをインストール
make install
```

## 開発

```bash
# 開発サーバーを起動
make dev

# もしくは
go run cmd/api/main.go
```

サーバーは http://localhost:8080 で起動します。

## API仕様

```bash
# ReDocでAPI仕様を確認
make redoc
```

ブラウザで http://localhost:3001 にアクセスしてAPI仕様を確認できます。

## 利用可能なコマンド

必要なコマンドは  make コマンドで全て確認することができます


