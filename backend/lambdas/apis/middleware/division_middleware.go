package middleware

import (
	"net/http"

	"apis/service"

	"github.com/labstack/echo/v5"
)

func CheckDivisionExist(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		id := ctx.Param("divisionID")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Division ID is required")
		}

		s := service.GetDivisionService()
		division, err := s.GetDivision(ctx, id)
		if err != nil {
			return err
		}

		ctx.Set("division", division)
		return next(ctx)
	}
}
