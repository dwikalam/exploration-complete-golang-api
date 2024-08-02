package transaction

import (
	"context"
	"database/sql"
)

func ExtractTxFromContext(ctx context.Context) (*sql.Tx, bool) {
	tx := ctx.Value(ctxTxKey)

	if v, ok := tx.(*sql.Tx); ok {
		return v, true
	}

	return nil, false
}

func CtxWithTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, ctxTxKey, tx)
}
