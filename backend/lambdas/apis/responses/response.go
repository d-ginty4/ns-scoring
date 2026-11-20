package responses

import (
	"github.com/labstack/echo/v5"
)

const requestIDContextKey = "requestID"

type Envelope[T any] struct {
	RequestID string    `json:"requestID"`
	Data      T         `json:"data"`
	Error     *APIError `json:"error"`
}

// APIError is the error object inside Envelope when a request fails.
type APIError struct {
	Message string   `json:"message" example:"validation failed"`
	Details []string `json:"details,omitempty" example:"'cursor' invalid cursor"`
}

func RequestID(c *echo.Context) string {
	if id, ok := c.Get(requestIDContextKey).(string); ok && id != "" {
		return id
	}
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

func Success(c *echo.Context, statusCode int, data any) error {
	id := RequestID(c)
	if id != "" {
		c.Response().Header().Set(echo.HeaderXRequestID, id)
	}
	data, err := omitLoadedRelationIDs(data)
	if err != nil {
		return err
	}
	return write(c, statusCode, Envelope[any]{
		RequestID: id,
		Data:      data,
		Error:     nil,
	})
}

func write[T any](c *echo.Context, statusCode int, body Envelope[T]) error {
	return c.JSON(statusCode, body)
}
