package database

import (
	"context"
	"database/sql"

	"github.com/stephenafamo/bob"
)

// WithTransaction はトランザクション内で任意の処理を実行します
// エラーが発生した場合は自動的にロールバック、成功時はコミットします
func WithTransaction(ctx context.Context, db *sql.DB, fn func(tx bob.Executor) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	executor := bob.NewTx(tx)
	if err := fn(executor); err != nil {
		return err
	}

	return tx.Commit()
}
