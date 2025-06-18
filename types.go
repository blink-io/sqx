package sqx

import (
	"context"

	"github.com/blink-io/sq"
)

const (
	DialectMySQL     = sq.DialectMySQL
	DialectPostgres  = sq.DialectPostgres
	DialectSQLite    = sq.DialectSQLite
	DialectSQLServer = sq.DialectSQLServer
)

var (
	MySQL     = sq.MySQL
	Postgres  = sq.Postgres
	SQLite    = sq.SQLite
	SQLServer = sq.SQLServer

	AlwaysTrueExpr = sq.Expr("1 = 1")
)

type (
	InsertQ func() sq.InsertQuery

	UpdateQ func() sq.UpdateQuery

	DeleteQ func() sq.DeleteQuery

	SelectQ[T any] func() (sq.SelectQuery, func(*sq.Row) T)

	Mapper[T sq.Table, M any, S any] interface {
		Table() T

		InsertT(ctx context.Context, ss ...S) func(*sq.Column)

		UpdateT(ctx context.Context, s S) func(*sq.Column)

		SelectT(ctx context.Context) func(*sq.Row) M
	}

	Executor[M any, S any] interface {
		Insert(ctx context.Context, db sq.DB, ss ...S) (sq.Result, error)

		Update(ctx context.Context, db sq.DB, where sq.Predicate, s S) (sq.Result, error)

		Delete(ctx context.Context, db sq.DB, where sq.Predicate) (sq.Result, error)

		One(ctx context.Context, db sq.DB, where sq.Predicate) (M, error)

		All(ctx context.Context, db sq.DB, where sq.Predicate) ([]M, error)
	}
)
