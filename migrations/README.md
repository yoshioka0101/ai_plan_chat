# ai_plan_chat

## マイグレーションの実行手順
### Database のセットアップ
1. データベースの起動
```bash
docker-compose up -d mysql
```

2. Atlas のマイグレーション実行
```bash
atlas migrate apply --env dev
```
