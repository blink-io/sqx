package sqx

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/bokwoon95/sq"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

var ctx = context.Background()

func TestDefaultDialect_1(t *testing.T) {
	dd := sq.DialectMySQL
	SetDefaultDialect(dd)

	fmt.Println("dialect before: ", *sq.DefaultDialect.Load())
	UnsetDefaultDialect()
	UnsetDefaultDialect()

	fmt.Println("dialect after: ", *sq.DefaultDialect.Load())
}

func TestInTxDB_1(t *testing.T) {
	dsn := "file:test.db?cache=shared&mode=memory"
	db, err := sql.Open("sqlite", dsn)
	require.NoError(t, err)

	txdb := InTx(db)

	err = txdb.RunInTx(ctx, nil, func(ctx context.Context, db sq.DB) error {
		q := sq.Queryf("select sqlite_version() as ver")

		ver, err := sq.FetchOne(sq.Log(db), q, func(r *sq.Row) string {
			return r.String("ver")
		})

		if err != nil {
			return nil
		}

		fmt.Println("sqlite version: ", ver)
		return nil
	})
	require.NoError(t, err)
}
