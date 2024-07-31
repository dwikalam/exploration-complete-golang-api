package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dwikalam/ecommerce-service/internal/app/types/interfaces"
)

type txKey string

var ctxTxKey = txKey("tx")

type SQLTransactionManager struct {
	db interfaces.DbAccessor
}

func NewManager(db interfaces.DbAccessor) (SQLTransactionManager, error) {
	if db == nil {
		return SQLTransactionManager{}, errors.New("interfaces.DbAccessor is nil")
	}

	return SQLTransactionManager{
		db: db,
	}, nil
}

func (t *SQLTransactionManager) Run(
	ctx context.Context,
	callback func(ctx context.Context) error,
) (rErr error) {
	var (
		tx    *sql.Tx
		txOpt sql.TxOptions = sql.TxOptions{}

		err error

		handleRecover = func() {
			rec := recover()

			if rec == nil {
				return
			}

			if e, ok := rec.(error); ok {
				rErr = e

				return
			}

			rErr = fmt.Errorf("%s", rec)
		}

		handleRollback = func() {
			if rErr != nil {
				rErr = errors.Join(rErr, tx.Rollback())
			}
		}
	)

	tx, err = t.db.Access().BeginTx(ctx, &txOpt)
	if err != nil {
		return err
	}

	defer handleRollback()
	defer handleRecover()

	if err = callback(ctxWithTx(ctx, tx)); err != nil {
		return err
	}

	return tx.Commit()
}

func ExtractTxFromContext(ctx context.Context) (*sql.Tx, bool) {
	tx := ctx.Value(ctxTxKey)

	if v, ok := tx.(*sql.Tx); ok {
		return v, true
	}

	return nil, false
}

func ctxWithTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, ctxTxKey, tx)
}
