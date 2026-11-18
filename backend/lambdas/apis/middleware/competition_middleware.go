package middleware

import (
	"net/http"

	"apis/service"

	"github.com/labstack/echo/v5"
)

func CheckCompetitionExist(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		id := ctx.Param("competitionID")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Competition ID is required")
		}

		s := service.GetCompetitionService()
		competition, err := s.GetCompetition(ctx, id)
		if err != nil {
			return err
		}

		ctx.Set("competition", competition)
		return next(ctx)
	}
}
