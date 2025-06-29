package sqx

import (
	"context"

	"github.com/blink-io/sq"
)

type Table[M any, S any] interface {
	sq.Table
	TableFunc[M, S]
}

type TableFunc[M any, S any] interface {
	ColumnMapper(ss ...S) sq.ColumnMapper

	RowMapper(context.Context, *sq.Row) M
}

type Executor[M any, S any] interface {
	Insert(ctx context.Context, db sq.DB, ss ...S) (sq.Result, error)

	Update(ctx context.Context, db sq.DB, where sq.Predicate, s S) (sq.Result, error)

	Delete(ctx context.Context, db sq.DB, where sq.Predicate) (sq.Result, error)

	One(ctx context.Context, db sq.DB, where sq.Predicate) (M, error)

	All(ctx context.Context, db sq.DB, where sq.Predicate) ([]M, error)
}

type executor[T Table[M, S], M any, S any] struct {
	t T
}

func NewExecutor[T Table[M, S], M any, S any](t T) Executor[M, S] {
	return executor[T, M, S]{t: t}
}

func (e executor[T, M, S]) Insert(ctx context.Context, db sq.DB, ss ...S) (sq.Result, error) {
	q := sq.InsertInto(e.t).
		ColumnValues(e.t.ColumnMapper(ss...))
	return sq.ExecContext(ctx, db, q)
}

func (e executor[T, M, S]) Update(ctx context.Context, db sq.DB, where sq.Predicate, s S) (sq.Result, error) {
	q := sq.Update(e.t).
		SetFunc(e.t.ColumnMapper(s)).
		Where(where)
	return sq.ExecContext(ctx, db, q)
}

func (e executor[T, M, S]) Delete(ctx context.Context, db sq.DB, where sq.Predicate) (sq.Result, error) {
	q := sq.DeleteFrom(e.t).
		Where(where)
	return sq.ExecContext(ctx, db, q)
}

func (e executor[T, M, S]) One(ctx context.Context, db sq.DB, where sq.Predicate) (M, error) {
	q := sq.From(e.t).
		Where(where).
		Limit(1)
	return sq.FetchOneContext[M](ctx, db, q, e.t.RowMapper)
}

func (e executor[T, M, S]) All(ctx context.Context, db sq.DB, where sq.Predicate) ([]M, error) {
	q := sq.From(e.t).
		Where(where)
	return sq.FetchAllContext[M](ctx, db, q, e.t.RowMapper)
}
