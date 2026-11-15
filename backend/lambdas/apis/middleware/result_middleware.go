package middleware

import (
	"errors"
	"net/http"

	"apis/exceptions"
	"apis/repository"

	"github.com/labstack/echo/v5"
)

func CheckResultExist(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		id := ctx.Param("resultID")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Result ID is required")
		}

		r := repository.GetResultRepository()
		result, err := r.GetResult(ctx, id)
		if err != nil {
			if errors.Is(err, &exceptions.ErrNotFound{}) {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		ctx.Set("result", result)
		return next(ctx)
	}
}
