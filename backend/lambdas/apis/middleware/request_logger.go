package middleware

import (
	"log/slog"
	"strings"

	"apis/logger"

	"github.com/labstack/echo/v5"
	echomw "github.com/labstack/echo/v5/middleware"
)

// AccessLog records every HTTP request after it completes: method, path, status,
// duration, requestID, and the caller when auth has identified them.
func AccessLog() echo.MiddlewareFunc {
	return echomw.RequestLoggerWithConfig(echomw.RequestLoggerConfig{
		Skipper: func(ctx *echo.Context) bool {
			return strings.HasPrefix(ctx.Request().URL.Path, "/swagger")
		},
		HandleError:  true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogMethod:    true,
		LogURIPath:   true,
		LogRoutePath: true,
		LogRequestID: true,
		LogStatus:    true,
		LogValuesFunc: func(ctx *echo.Context, v echomw.RequestLoggerValues) error {
			attrs := []any{
				slog.String("method", v.Method),
				slog.String("path", v.URIPath),
				slog.String("route", v.RoutePath),
				slog.Int("status", v.Status),
				slog.Int64("durationMs", v.Latency.Milliseconds()),
				slog.String("requestID", v.RequestID),
				slog.String("remoteIp", v.RemoteIP),
			}

			if userID, ok := ctx.Get("userID").(string); ok && userID != "" {
				attrs = append(attrs, slog.String("userID", userID))
			}
			// if user, ok := ctx.Get("requestUser").(*models.User); ok && user != nil {
			// 	if user.Name != "" {
			// 		attrs = append(attrs, slog.String("userName", user.Name))
			// 	}
			// 	if user.IsAdmin {
			// 		attrs = append(attrs, slog.Bool("isAdmin", true))
			// 	}
			// }
			if v.Error != nil {
				attrs = append(attrs, logger.Err(v.Error))
			}

			level := slog.LevelInfo
			switch {
			case v.Status >= 500:
				level = slog.LevelError
			case v.Status >= 400:
				level = slog.LevelWarn
			}

			logger.Logger().Log(ctx.Request().Context(), level, "request", attrs...)
			return nil
		},
	})
}

// RequestID assigns a request id (from X-Request-ID or generated) and binds a
// request-scoped logger that downstream code should use via logger.From(c).
func RequestID() echo.MiddlewareFunc {
	return echomw.RequestIDWithConfig(echomw.RequestIDConfig{
		RequestIDHandler: func(c *echo.Context, requestID string) {
			logger.BindRequest(c, requestID)
		},
	})
}
