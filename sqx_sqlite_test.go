package sqx

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/blink-io/sq"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

var ctx = context.Background()

func MustGetSQLite() *sql.DB {
	dsn := "file:test.db?cache=shared&mode=memory"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		panic(err)
	}
	return db
}

func querySQLiteVersion(ctx context.Context, db sq.DB) error {
	q := sq.Queryf("select sqlite_version() as ver")

	ver, err := sq.FetchOne(sq.Log(db), q, func(ctx context.Context, r *sq.Row) string {
		return r.String("ver")
	})

	if err != nil {
		return nil
	}

	fmt.Println("sqlite version: ", ver)

	return nil
}

func TestDefaultDialect_1(t *testing.T) {
	dd := sq.DialectMySQL
	SetDefaultDialect(dd)

	fmt.Println("dialect before: ", *sq.DefaultDialect.Load())
	RestoreDefaultDialect()

	fmt.Println("dialect after: ", *sq.DefaultDialect.Load())
}

func TestInTxDB_1(t *testing.T) {
	db := MustGetSQLite()

	txdb := InTx(db)

	err := txdb.RunInTx(ctx, nil, func(ctx context.Context, db sq.DB) error {
		return querySQLiteVersion(ctx, db)
	})
	require.NoError(t, err)
}

func TestHookDB_1(t *testing.T) {
	db := MustGetSQLite()

	hdb := Hooks(sq.Log(db), logHook{
		logf: func(ctx context.Context, format string, args ...any) {
			fmt.Printf(format, args...)
		},
	})

	err := querySQLiteVersion(ctx, hdb)
	require.NoError(t, err)
}

func TestExecutor_1(t *testing.T) {

}
