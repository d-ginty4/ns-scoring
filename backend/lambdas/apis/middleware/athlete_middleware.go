package middleware

import (
	"net/http"

	"apis/service"

	"github.com/labstack/echo/v5"
)

func CheckAthleteExist(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		id := ctx.Param("athleteID")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Athlete ID is required")
		}

		s := service.GetAthleteService()
		athlete, err := s.GetAthlete(ctx, ctx.Param("competitionID"), id)
		if err != nil {
			return err
		}

		ctx.Set("athlete", athlete)
		return next(ctx)
	}
}
