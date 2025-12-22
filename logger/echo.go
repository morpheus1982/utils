package logger

import (
	"sync/atomic"

	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap/zapcore"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var concurrencyCount int64

func Middleware() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		BeforeNextFunc: func(c echo.Context) {
			atomic.AddInt64(&concurrencyCount, 1)
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			count := atomic.AddInt64(&concurrencyCount, -1)
			fields := []zapcore.Field{
				zap.String("url", v.URI),
				zap.String("method", v.Method),
				zap.Int("status", v.Status),
				zap.Duration("latency", v.Latency),
				zap.String("host", v.Host),
				zap.String("referer", v.Referer),
				zap.String("user_agent", v.UserAgent),
				zap.String("remote_addr", c.Request().RemoteAddr),
				zap.String("request_id", v.RequestID),
				zap.String("content_length", v.ContentLength),
				zap.Int64("response_size", v.ResponseSize),
				zap.Int64("concurrency_count", count),
			}
			if val := c.Get("fields"); val != nil {
				fs := val.([]string)
				for _, f := range fs {
					fields = append(fields, zap.Any(f, c.Get(f)))
				}
			}
			if v.Error != nil {
				fields = append(fields, zap.Error(v.Error))
			}
			access(fields...)
			return nil
		},
		LogLatency:       true,
		LogMethod:        true,
		LogURI:           true,
		LogURIPath:       true,
		LogRoutePath:     true,
		LogRequestID:     true,
		LogReferer:       true,
		LogUserAgent:     true,
		LogStatus:        true,
		LogError:         true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogRemoteIP:      true,
	})
}
