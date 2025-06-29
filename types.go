package sqx

import (
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
)
