package sqx

import (
	"context"

	"github.com/bokwoon95/sq"
)

var _ Executor[any, any] = executor[sq.Table, any, any]{}

type executor[T sq.Table, M any, S any] struct {
	m Mapper[T, M, S]
}

func NewExecutor[T sq.Table, M any, S any](m Mapper[T, M, S]) Executor[M, S] {
	return executor[T, M, S]{m: m}
}

func (e executor[T, M, S]) Insert(ctx context.Context, db sq.DB, ss ...S) (sq.Result, error) {
	q := sq.InsertInto(e.m.Table()).
		ColumnValues(e.m.InsertT(ctx, ss...))
	return sq.ExecContext(ctx, db, q)
}

func (e executor[T, M, S]) Update(ctx context.Context, db sq.DB, where sq.Predicate, s S) (sq.Result, error) {
	q := sq.Update(e.m.Table()).
		SetFunc(e.m.UpdateT(ctx, s)).
		Where(where)
	return sq.ExecContext(ctx, db, q)
}

func (e executor[T, M, S]) Delete(ctx context.Context, db sq.DB, where sq.Predicate) (sq.Result, error) {
	q := sq.DeleteFrom(e.m.Table()).
		Where(where)
	return sq.ExecContext(ctx, db, q)
}

func (e executor[T, M, S]) One(ctx context.Context, db sq.DB, where sq.Predicate) (M, error) {
	q := sq.From(e.m.Table()).
		Where(where).
		Limit(1)
	return sq.FetchOneContext[M](ctx, db, q, e.m.SelectT(ctx))
}

func (e executor[T, M, S]) All(ctx context.Context, db sq.DB, where sq.Predicate) ([]M, error) {
	q := sq.From(e.m.Table()).
		Where(where)
	return sq.FetchAllContext[M](ctx, db, q, e.m.SelectT(ctx))
}
