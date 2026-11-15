package middleware

import (
	"errors"
	"net/http"

	"apis/exceptions"
	"apis/repository"

	"github.com/labstack/echo/v5"
)

func CheckEventExist(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		id := ctx.Param("eventID")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Event ID is required")
		}

		r := repository.GetEventRepository()
		event, err := r.GetEvent(ctx, id)
		if err != nil {
			if errors.Is(err, &exceptions.ErrNotFound{}) {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		ctx.Set("event", event)
		return next(ctx)
	}
}
