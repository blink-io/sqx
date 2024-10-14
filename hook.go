package sqx

import (
	"context"
	"database/sql"
	"time"

	"github.com/bokwoon95/sq"
)

type (
	QueryEvent struct {
		DB *DB

		Query     string
		Args      []any
		StartTime time.Time
		Result    sql.Result
		Err       error

		Stash map[any]any
	}

	QueryHook interface {
		BeforeQuery(context.Context, *QueryEvent) context.Context
		AfterQuery(context.Context, *QueryEvent)
	}

	hookDB struct {
		sq.DB
		hooks []QueryHook
	}
)

func Hooks(db sq.DB, hooks ...QueryHook) interface {
	sq.DB
} {
	return hookDB{DB: db, hooks: hooks}
}

func (db hookDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	ctx, event := db.beforeQuery(ctx, query, args...)
	rows, err := db.DB.QueryContext(ctx, query, args...)
	db.afterQuery(ctx, event, nil, err)
	return rows, err
}

func (db hookDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	ctx, event := db.beforeQuery(ctx, query, args...)
	res, err := db.DB.ExecContext(ctx, query, args...)
	db.afterQuery(ctx, event, res, err)
	return res, err
}

func (db hookDB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return db.DB.PrepareContext(ctx, query)
}

func (db hookDB) beforeQuery(ctx context.Context, query string, args ...any) (context.Context, *QueryEvent) {
	if len(db.hooks) == 0 {
		return ctx, nil
	}

	event := &QueryEvent{
		Query: query,
		Args:  args,

		StartTime: time.Now(),
	}

	for _, hook := range db.hooks {
		ctx = hook.BeforeQuery(ctx, event)
	}

	return ctx, event
}

func (db hookDB) afterQuery(
	ctx context.Context,
	event *QueryEvent,
	res sql.Result,
	err error,
) {
	switch err {
	case nil, sql.ErrNoRows:
		// nothing
	default:
		//atomic.AddUint32(&db.stats.Errors, 1)
	}

	if event == nil {
		return
	}

	event.Result = res
	event.Err = err

	db.afterQueryFromIndex(ctx, event, len(db.hooks)-1)
}

func (db hookDB) afterQueryFromIndex(ctx context.Context, event *QueryEvent, hookIndex int) {
	for ; hookIndex >= 0; hookIndex-- {
		db.hooks[hookIndex].AfterQuery(ctx, event)
	}
}
