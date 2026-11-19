package middleware

import (
	"net/http"

	"apis/service"

	"github.com/labstack/echo/v5"
)

func CheckResultExist(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		eventID := ctx.Param("eventID")
		if eventID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Event ID is required")
		}

		s := service.GetResultService()
		result, err := s.GetResult(ctx, ctx.Param("athleteID"), eventID)
		if err != nil {
			return err
		}

		ctx.Set("result", result)
		return next(ctx)
	}
}
