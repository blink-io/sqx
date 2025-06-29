package sqx

import (
	"log/slog"
	"strings"

	"github.com/blink-io/sq"
)

type (
	DB        = sq.DB
	DBX       = TxDB
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
		slog.Warn("unsupported dialect")
	}
}

// RestoreDefaultDialect restores default dialect as unset.
func RestoreDefaultDialect() {
	sq.DefaultDialect.Store(nil)
}
