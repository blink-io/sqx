package sqx

import (
	"context"
	"github.com/blink-io/sq"
)

type (
	Mapper[T sq.Table, M any, S any] interface {
		Table() T

		InsertT(ctx context.Context, ss ...S) sq.ColumnMapper

		UpdateT(ctx context.Context, s S) sq.ColumnMapper

		SelectT(ctx context.Context) sq.RowMapper[M]
	}

	MapperTable[M any, S any] interface {
		sq.Table
		TableFunc[M, S]
	}

	TableFunc[M any, S any] interface {
		ColumnMapper(ss ...S) sq.ColumnMapper

		RowMapper(context.Context, *sq.Row) M

		RowMapperFunc() sq.RowMapper[M]
	}

	mapper[T MapperTable[M, S], M any, S any] struct {
		t T
	}
)

func (m mapper[T, M, S]) Table() T {
	return m.t
}

func (m mapper[T, M, S]) InsertT(ctx context.Context, ss ...S) sq.ColumnMapper {
	return m.t.ColumnMapper(ss...)
}

func (m mapper[T, M, S]) UpdateT(ctx context.Context, s S) sq.ColumnMapper {
	return m.t.ColumnMapper(s)
}

func (m mapper[T, M, S]) SelectT(ctx context.Context) sq.RowMapper[M] {
	return m.t.RowMapper
}

func NewMapper[T MapperTable[M, S], M any, S any](t T) Mapper[T, M, S] {
	return mapper[T, M, S]{t: t}
}
