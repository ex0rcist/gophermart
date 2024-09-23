package tracer

import (
	"context"
	"fmt"
	"time"

	"github.com/ex0rcist/gophermart/internal/logging"
	"github.com/jackc/pgx/v5"
)

const (
	queryStartTimeKey contextKey = "queryStartTime"
	sqlQueryKey       contextKey = "sqlQuery"
)

type contextKey string

type DBQueryTracer struct {
	logFunc func(ctx context.Context, sqlQuery any, duration time.Duration)
}

func NewDBQueryTracer() *DBQueryTracer {
	return &DBQueryTracer{logFunc: defaultLogFunc}
}

func (tracer *DBQueryTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	ctx = context.WithValue(ctx, queryStartTimeKey, time.Now())
	ctx = context.WithValue(ctx, sqlQueryKey, data.SQL)

	return ctx
}

func (tracer *DBQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	startTime, ok := ctx.Value(queryStartTimeKey).(time.Time)
	if !ok {
		startTime = time.Now()
	}

	duration := time.Since(startTime)
	sqlQuery := ctx.Value(sqlQueryKey)

	tracer.logFunc(ctx, sqlQuery, duration)
}

func defaultLogFunc(ctx context.Context, sqlQuery any, duration time.Duration) {
	logging.LogInfoCtx(ctx, fmt.Sprintf("pgx: queried \"%s\" in %v", sqlQuery, duration))
}
