package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/ex0rcist/gophermart/internal/logging"
	"github.com/jackc/pgx/v5"
)

type dbQueryTracer struct {
}

type contextKey string

const (
	queryStartTimeKey contextKey = "queryStartTime"
	sqlQueryKey       contextKey = "sqlQuery"
)

func (tracer *dbQueryTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	ctx = context.WithValue(ctx, queryStartTimeKey, time.Now())
	ctx = context.WithValue(ctx, sqlQueryKey, data.SQL)

	return ctx
}

func (tracer *dbQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	startTime, ok := ctx.Value(queryStartTimeKey).(time.Time)
	if !ok {
		startTime = time.Now()
	}

	duration := time.Since(startTime)
	sqlQuery := ctx.Value(sqlQueryKey)

	logging.LogDebugCtx(ctx, fmt.Sprintf("pgx: queried \"%s\" in %v", sqlQuery, duration))
}
