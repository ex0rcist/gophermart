package accrual

import (
	"context"

	"github.com/ex0rcist/gophermart/internal/utils"
	"github.com/rs/zerolog/log"
)

func setupCtxWithRID(ctx context.Context) context.Context {
	logger := log.Logger.With().Ctx(ctx).Str("rid", utils.GenerateRequestID()).Logger()
	return logger.WithContext(ctx)
}
