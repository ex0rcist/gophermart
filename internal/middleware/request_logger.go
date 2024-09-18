package middleware

import (
	"net/http"
	"time"

	"github.com/ex0rcist/gophermart/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func RequestsLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := findOrCreateRequestID(c.Request)

		// setup child logger for middleware
		logger := log.Logger.With().
			Str("rid", requestID).
			Logger()

		// log started
		logger.Info().
			Str("method", c.Request.Method).
			Str("url", c.Request.URL.String()).
			Str("remote-addr", c.Request.RemoteAddr).
			Msg("Started")

		logger.Debug().
			Msgf("request: %s", utils.HeadersToStr(c.Request.Header))

		c.Writer.Header().Set("X-Request-Id", requestID)

		ctx := logger.WithContext(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)

		// execute
		c.Next()

		if headers := c.Writer.Header(); len(headers) > 0 {
			logger.Debug().
				Msgf("response: %s", utils.HeadersToStr(headers))
		}

		// log completed
		logger.Info().
			Float64("elapsed", time.Since(start).Seconds()).
			Int("status", c.Writer.Status()).
			Msg("Completed")
	}
}

func findOrCreateRequestID(r *http.Request) string {
	requestID := r.Header.Get("X-Request-Id")

	if requestID == "" {
		requestID = utils.GenerateRequestID()
	}

	return requestID
}
