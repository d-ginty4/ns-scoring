package middleware

import (
	"errors"
	"net/http"

	"apis/exceptions"
	"apis/repository"

	"github.com/labstack/echo/v5"
)

func CheckDivisionExist(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		id := ctx.Param("divisionID")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Division ID is required")
		}

		r := repository.GetDivisionRepository()
		division, err := r.GetDivision(ctx, id)
		if err != nil {
			if errors.Is(err, &exceptions.ErrNotFound{}) {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		ctx.Set("division", division)
		return next(ctx)
	}
}
