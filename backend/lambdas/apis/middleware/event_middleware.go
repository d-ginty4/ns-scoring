package middleware

import (
	"net/http"

	"apis/service"

	"github.com/labstack/echo/v5"
)

func CheckEventExist(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		id := ctx.Param("eventID")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Event ID is required")
		}

		s := service.GetEventService()
		event, err := s.GetEvent(ctx, id)
		if err != nil {
			return err
		}

		ctx.Set("event", event)
		return next(ctx)
	}
}
