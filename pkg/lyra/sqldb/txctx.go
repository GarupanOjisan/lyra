package sqldb

import (
    "context"
    "database/sql"
)

type txKey struct{}

func WithTx(ctx context.Context, tx *sql.Tx) context.Context { return context.WithValue(ctx, txKey{}, tx) }
func TxFrom(ctx context.Context) *sql.Tx {
    if v := ctx.Value(txKey{}); v != nil {
        if tx, ok := v.(*sql.Tx); ok { return tx }
    }
    return nil
}

