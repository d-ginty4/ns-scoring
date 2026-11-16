package routes

import (
	"net/http"

	"apis/models"
	"apis/responses"
	"apis/service"
	"apis/utils"

	"github.com/labstack/echo/v5"
)

type RegionRoutes struct {
	QueryRegions func(c *echo.Context) error
	PostRegion   func(c *echo.Context) error
	GetRegion    func(c *echo.Context) error
	PutRegion    func(c *echo.Context) error
	DeleteRegion func(c *echo.Context) error
}

var regionRoutes *RegionRoutes

func getRegionRoutes() *RegionRoutes {
	if regionRoutes == nil {
		return newRegionRoutes()
	}
	return regionRoutes
}

func newRegionRoutes() *RegionRoutes {
	routes := &RegionRoutes{
		QueryRegions: queryRegions,
		PostRegion:   postRegion,
		GetRegion:    getRegion,
		PutRegion:    putRegion,
		DeleteRegion: deleteRegion,
	}

	regionRoutes = routes
	return routes
}

func queryRegions(ctx *echo.Context) error {
	service := service.GetRegionService()

	page, err := service.QueryRegions(ctx)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, page)
}

func postRegion(ctx *echo.Context) error {
	service := service.GetRegionService()

	var p models.RegionRequest
	if err := utils.BindAndValidate(ctx, &p); err != nil {
		return err
	}

	region, err := service.CreateRegion(ctx, &p)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, region)
}

func getRegion(ctx *echo.Context) error {
	region, ok := ctx.Get("region").(*models.Region)
	if ok && region != nil {
		return responses.Success(ctx, http.StatusOK, region)
	}

	regionService := service.GetRegionService()
	id := ctx.Param("regionID")

	region, err := regionService.GetRegion(ctx, id)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, region)
}

func putRegion(ctx *echo.Context) error {
	service := service.GetRegionService()

	var p models.RegionRequest
	if err := utils.BindAndValidate(ctx, &p); err != nil {
		return err
	}

	id := ctx.Param("regionID")
	existingRegion, _ := ctx.Get("region").(*models.Region)
	region, err := service.UpdateRegion(ctx, id, &p, existingRegion)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, region)
}

func deleteRegion(ctx *echo.Context) error {
	service := service.GetRegionService()

	id := ctx.Param("regionID")
	if err := service.DeleteRegion(ctx, id); err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, nil)
}
