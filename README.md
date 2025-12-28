# ai_plan_chat

AIを介してタスクを作成するアプリ。タスク管理が面倒なときに、自然文で入力するだけでタスク化できます。

## セットアップ

### 必要環境
- Go 1.24
- Node.js 18+
- MySQL 8.4（ローカルは `migrations/docker-compose.yml` を利用）
- Atlas（マイグレーション用）

### 環境変数

バックエンド（`backend/.env` を作成して設定）:

```env
# Database (either DB_DSN or these three)
DB_USER=root
DB_PASSWORD=password
DB_NAME=ai_chat_task
DB_HOST=127.0.0.1
DB_PORT=3306

# Auth
JWT_SECRET=your-32+chars-secret
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback

# AI
GEMINI_API_KEY=your-gemini-api-key
GEMINI_MODEL=gemini-2.5-flash-lite
```

フロントエンド（`frontend/.env` を作成して設定）:

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_GOOGLE_AUTH_URL=http://localhost:8080/auth/google/callback
VITE_ENV=development
```

### 開発起動

```bash
# DB起動
cd migrations
docker-compose up -d mysql

# マイグレーション実行
atlas migrate apply --env dev
```

```bash
# Backend
cd backend
make dev
```

```bash
# Frontend
cd frontend
npm install
npm run dev
```

## API 設計書

API 設計書の定義場所は `backend/internal/api/openapi/openapi.yaml` に記載しています。

### オンライン版
GitHub Pages で公開されたAPI設計書: https://yoshioka0101.github.io/ai_plan_chat/

### ローカル確認
ローカルで確認したい場合は `backend` ディレクトリで以下のコマンドを実行:

```bash
# ReDoc でOpenAPI仕様を確認
make redoc
```

## ディレクトリ構造

```
.
├── backend/       # Go APIサーバー
│   ├── cmd/       # エントリーポイント
│   ├── internal/  # アプリ本体（API/Usecase/Repoなど）
│   └── Makefile   # 開発コマンド
├── frontend/      # React + Vite フロントエンド
│   ├── src/       # UI/ページ/サービス
│   └── public/    # 静的ファイル
├── migrations/    # AtlasマイグレーションとDB起動設定
│   ├── atlas.hcl  # Atlas設定
│   └── docker-compose.yml
└── schemas/       # SQLスキーマ定義
    └── schema.sql
```

## 技術スタック
- Backend: Go, Gin, BOB, oapi-codegen
- Frontend: React, TypeScript, Vite, Axios, React Router
- DB: MySQL 8.4, Atlas
- AI: Google Gemini API

## 開発のルールについて
- issuesに仕様書を記載する

### ブランチ戦略について
- mainブランチ
- リリースブランチ:v1.x.x
- 開発ブランチ:issues/xx

## ER図
- {ここにER図の全体像のリンクをおいておく}

## 自動デプロイについて
- 未設定
