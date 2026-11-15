package middleware

import (
	"errors"
	"net/http"

	"apis/exceptions"
	"apis/repository"

	"github.com/labstack/echo/v5"
)

func CheckAthleteExist(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		id := ctx.Param("athleteID")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Athlete ID is required")
		}

		r := repository.GetAthleteRepository()
		athlete, err := r.GetAthlete(ctx, id)
		if err != nil {
			if errors.Is(err, &exceptions.ErrNotFound{}) {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		ctx.Set("athlete", athlete)
		return next(ctx)
	}
}
