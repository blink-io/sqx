package sqx

import (
	"log/slog"
	"strings"

	"github.com/bokwoon95/sq"
)

type (
	DB        = sq.DB
	Predicate = sq.Predicate
	Query     = sq.Query
	Row       = sq.Row
	SQLWriter = sq.SQLWriter

	JSONMap map[string]any
)

func SetDefaultDialect(dialect string) {
	switch dialect := strings.ToLower(dialect); dialect {
	case sq.DialectPostgres,
		sq.DialectSQLite,
		sq.DialectSQLServer,
		sq.DialectMySQL:
		sq.DefaultDialect.Store(&dialect)
	default:
		slog.Warn("")
	}
}

func UnsetDefaultDialect() {
	sq.DefaultDialect.Store(nil)
}
