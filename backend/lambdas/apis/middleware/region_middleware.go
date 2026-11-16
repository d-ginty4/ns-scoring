package middleware

import (
	"net/http"

	"apis/service"

	"github.com/labstack/echo/v5"
)

func CheckRegionExist(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		id := ctx.Param("regionID")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Region ID is required")
		}

		s := service.GetRegionService()
		region, err := s.GetRegion(ctx, id)
		if err != nil {
			return err
		}

		ctx.Set("region", region)
		return next(ctx)
	}
}
