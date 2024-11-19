package tracer

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestTraceQueryStart(t *testing.T) {
	tracer := NewDBQueryTracer()
	ctx := context.Background()

	queryData := pgx.TraceQueryStartData{
		SQL: "SELECT * FROM users",
	}

	ctx = tracer.TraceQueryStart(ctx, nil, queryData)

	startTime, ok := ctx.Value(queryStartTimeKey).(time.Time)
	assert.True(t, ok, "Expected queryStartTimeKey to be set in context")
	assert.WithinDuration(t, time.Now(), startTime, time.Second, "Expected start time to be close to current time")

	sqlQuery, ok := ctx.Value(sqlQueryKey).(string)
	assert.True(t, ok, "Expected sqlQueryKey to be set in context")
	assert.Equal(t, "SELECT * FROM users", sqlQuery, "Expected SQL query to be set correctly")
}

func TestTraceQueryEndWithMockLog(t *testing.T) {
	var loggedSQL any
	var loggedDuration time.Duration
	mockLog := func(ctx context.Context, sqlQuery any, duration time.Duration) {
		loggedSQL = sqlQuery
		loggedDuration = duration
	}

	tracer := NewDBQueryTracer()
	tracer.logFunc = mockLog

	ctx := context.Background()

	startTime := time.Now().Add(-3 * time.Second)
	ctx = context.WithValue(ctx, queryStartTimeKey, startTime)
	ctx = context.WithValue(ctx, sqlQueryKey, "SELECT * FROM users")

	queryEndData := pgx.TraceQueryEndData{}
	tracer.TraceQueryEnd(ctx, nil, queryEndData)

	assert.Equal(t, "SELECT * FROM users", loggedSQL, "Expected SQL query to be logged")
	assert.GreaterOrEqual(t, loggedDuration.Seconds(), 3.0, "Expected duration to be greater or equal to 3 seconds")
}
