# AI Plan Chat の Backend開発
Go言語で構築されたREST APIサーバーです。

### ディレクトリ構成

```
backend/
├── cmd/                      # アプリケーションエントリーポイント
│   ├── api/                 # APIサーバー起動
│   └── migrate/             # DBマイグレーション実行
│
├── config/                   # 設定管理（環境変数、設定ファイル読み込み）
│
├── internal/
│   ├── entity/              # ドメインエンティティ（データ構造のみ）
│   ├── repository/          # Repositoryインターフェース定義
│   ├── infrastructure/      # Repository実装（DB依存コード）
│   ├── usecase/             # ユースケース層（ビジネスロジック）
│   │
│   ├── http/                # HTTP層
│   │   ├── handler/        # HTTPリクエスト処理
│   │   ├── presenter/      # レスポンス変換（UseCaseの結果をHTTPレスポンスに変換）
│   │   └── routes.go       # ルーティング設定
│   │
│   ├── middleware/          # HTTPミドルウェア
│   ├── errors/              # カスタムエラー型
│   │
│   ├── api/openapi/         # OpenAPI定義ファイル
│   │   ├── openapi.yaml    # API仕様（分割版）
│   │   ├── components/     # 再利用可能なコンポーネント
│   │   └── paths/          # エンドポイント定義
│   │
│   └── gen/                 # 自動生成コード
│       ├── api/            # OpenAPI生成コード (oapi-codegen)
│       └── db/             # BOB生成コード (bobgen)
│
├── Makefile                 # ビルド・開発コマンド
├── bobgen.yaml              # BOB型生成設定
└── go.mod                   # Go依存関係
```

## 技術スタック

- フレームワーク: [Gin](https://github.com/gin-gonic/gin)
- クエリビルダー: [BOB](https://github.com/stephenafamo/bob) (型安全なSQLクエリビルダー)
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

## Dockerでの実行（ECS/ECR向け）
Dockerが使える環境で、バックエンドのみコンテナ起動する場合。

```bash
cd backend
# ARM向けイメージをビルド（ポートは8080）
docker buildx build --platform linux/arm64 -t hubplanner-backend:local -f Dockerfile .

# ローカル実行
docker run --rm -p 8080:8080 hubplanner-backend:local
# → http://localhost:8080/health が返ればOK
```

ECRにpushする場合は、リポジトリURIでタグ付けしてから `docker push` してください。

## 利用可能なコマンド

必要なコマンドは  make コマンドで全て確認することができます

