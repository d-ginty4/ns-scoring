package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/labstack/echo/v5"
)

const requestLoggerKey = "requestLogger"

var logger *slog.Logger

func init() {
	opts := &slog.HandlerOptions{Level: levelFromEnv()}
	logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

// Logger returns the process-wide structured logger.
func Logger() *slog.Logger {
	if logger == nil {
		return slog.Default()
	}
	return logger
}

// From returns the request-scoped logger when one is bound, otherwise the default logger.
func From(c *echo.Context) *slog.Logger {
	if c == nil {
		return Logger()
	}
	return c.Logger()
}

// BindRequest attaches a request-scoped logger with requestID, method, and path.
func BindRequest(c *echo.Context, requestID string) {
	req := c.Request()
	attrs := []any{slog.String("requestID", requestID)}
	if req != nil {
		attrs = append(attrs,
			slog.String("method", req.Method),
			slog.String("path", req.URL.Path),
		)
	}

	requestLogger := Logger().With(attrs...)
	c.Set("requestID", requestID)
	c.Set(requestLoggerKey, requestLogger)
	c.SetLogger(requestLogger)
}

// SetActor adds the known caller to the request-scoped logger.
func SetActor(c *echo.Context, userID, userName string, isAdmin bool) {
	if userID == "" {
		return
	}

	c.Set("userID", userID)

	attrs := []any{slog.String("userID", userID)}
	if userName != "" {
		attrs = append(attrs, slog.String("userName", userName))
	}
	if isAdmin {
		attrs = append(attrs, slog.Bool("isAdmin", true))
	}

	c.SetLogger(baseRequestLogger(c).With(attrs...))
}

// Action logs a completed business action. Extra attrs should be key/value pairs
// (for example "postID", id). The request logger already includes requestID and actor.
func Action(c *echo.Context, action string, attrs ...any) {
	From(c).Info("action", append([]any{slog.String("action", action)}, attrs...)...)
}

// Err returns a structured error attribute named "error".
func Err(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}
	return slog.Any("error", err)
}

func baseRequestLogger(c *echo.Context) *slog.Logger {
	if l, ok := c.Get(requestLoggerKey).(*slog.Logger); ok && l != nil {
		return l
	}
	return From(c)
}

func levelFromEnv() slog.Level {
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		switch strings.ToLower(v) {
		case "debug":
			return slog.LevelDebug
		case "warn", "warning":
			return slog.LevelWarn
		case "error":
			return slog.LevelError
		default:
			return slog.LevelInfo
		}
	}

	if os.Getenv("APP_ENV") == "local" {
		return slog.LevelDebug
	}

	return slog.LevelInfo
}
