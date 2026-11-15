package middleware

import (
	"errors"
	"net/http"

	"apis/exceptions"
	"apis/repository"

	"github.com/labstack/echo/v5"
)

func CheckCompetitionExist(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		id := ctx.Param("competitionID")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Competition ID is required")
		}

		r := repository.GetCompetitionRepository()
		competition, err := r.GetCompetition(ctx, id)
		if err != nil {
			if errors.Is(err, &exceptions.ErrNotFound{}) {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		ctx.Set("competition", competition)
		return next(ctx)
	}
}
