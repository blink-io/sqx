package slog

import (
	"context"
	"log/slog"

	"github.com/blink-io/sq"
)

type logger struct {
	ctx context.Context
	sl  *slog.Logger
	lv  slog.Leveler
	cfg sq.LoggerConfig
}

var _ sq.Logger = (*logger)(nil)

func New(sl *slog.Logger, lv slog.Leveler, cfg sq.LoggerConfig) sq.Logger {
	return &logger{
		ctx: context.Background(),
		sl:  sl,
		lv:  lv,
		cfg: cfg,
	}
}

func (l *logger) LogSettings(ctx context.Context, settings *sq.LogSettings) {
	settings.LogAsynchronously = l.cfg.LogAsynchronously
	settings.IncludeTime = l.cfg.ShowTimeTaken
	settings.IncludeCaller = l.cfg.ShowCaller
	settings.IncludeResults = l.cfg.ShowResults
}

func (l *logger) LogQuery(ctx context.Context, stats sq.QueryStats) {
	attrs := []slog.Attr{
		slog.String("dialect", stats.Dialect),
		slog.String("query", stats.Query),
		slog.Time("started_at", stats.StartedAt),
		slog.String("results", stats.Results),
	}
	if l.cfg.ShowTimeTaken {
		attrs = append(attrs, slog.Duration("time_taken", stats.TimeTaken))
	}
	if l.cfg.ShowCaller {
		attrs = append(attrs,
			slog.String("caller_file", stats.CallerFile),
			slog.Int("caller_line", stats.CallerLine),
			slog.String("caller_function", stats.CallerFunction),
		)
	}
	l.sl.LogAttrs(l.ctx, l.lv.Level(), "Query stats", attrs...)
}
