package sqx

import (
	"context"
	"database/sql"

	"github.com/bokwoon95/sq"
)

type (
	Txer interface {
		Begin() (*sql.Tx, error)
		BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	}

	TxDB interface {
		sq.DB
		Txer
	}

	runInTxer interface {
		RunInTx(context.Context, *sql.TxOptions, func(context.Context, sq.DB) error) error
	}

	inTxDB struct {
		TxDB
	}
)

func (db inTxDB) RunInTx(ctx context.Context, opts *sql.TxOptions, fn func(context.Context, sq.DB) error) error {
	return RunInTx(ctx, db, opts, fn)
}

func InTx(db TxDB) interface {
	TxDB
	runInTxer
} {
	return inTxDB{TxDB: db}
}

func RunInTx(ctx context.Context,
	db Txer, opts *sql.TxOptions, fn func(context.Context, sq.DB) error) error {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	var done bool

	defer func() {
		if !done {
			_ = tx.Rollback()
		}
	}()

	if err := fn(ctx, tx); err != nil {
		return err
	}

	done = true
	return tx.Commit()
}
