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

// queryRegions lists regions.
//
//	@Summary		List regions
//	@Tags			Regions
//	@Produce		json
//	@Success		200	{object}	responses.Envelope[[]models.Region]
//	@Failure		500	{object}	responses.Envelope[any]
//	@Router			/v1/regions [get]
func queryRegions(ctx *echo.Context) error {
	service := service.GetRegionService()

	page, err := service.QueryRegions(ctx)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, page)
}

// postRegion creates a region.
//
//	@Summary		Create region
//	@Tags			Regions
//	@Accept			json
//	@Produce		json
//	@Param			region	body		models.RegionRequest	true	"Region"
//	@Success		200		{object}	responses.Envelope[models.Region]
//	@Failure		400		{object}	responses.Envelope[any]
//	@Failure		500		{object}	responses.Envelope[any]
//	@Router			/v1/regions [post]
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

// getRegion returns a region.
//
//	@Summary		Get region
//	@Tags			Regions
//	@Produce		json
//	@Param			regionID	path		integer	true	"Region ID"
//	@Success		200			{object}	responses.Envelope[models.Region]
//	@Failure		404			{object}	responses.Envelope[any]
//	@Failure		500			{object}	responses.Envelope[any]
//	@Router			/v1/regions/{regionID} [get]
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

// putRegion creates or updates a region.
//
//	@Summary		Update region
//	@Tags			Regions
//	@Accept			json
//	@Produce		json
//	@Param			regionID	path		integer					true	"Region ID"
//	@Param			region		body		models.RegionRequest	true	"Region"
//	@Success		200			{object}	responses.Envelope[models.Region]
//	@Failure		400			{object}	responses.Envelope[any]
//	@Failure		404			{object}	responses.Envelope[any]
//	@Failure		500			{object}	responses.Envelope[any]
//	@Router			/v1/regions/{regionID} [put]
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

// deleteRegion deletes a region.
//
//	@Summary		Delete region
//	@Tags			Regions
//	@Produce		json
//	@Param			regionID	path		integer	true	"Region ID"
//	@Success		200			{object}	responses.Envelope[models.Region]
//	@Failure		404			{object}	responses.Envelope[any]
//	@Failure		500			{object}	responses.Envelope[any]
//	@Router			/v1/regions/{regionID} [delete]
func deleteRegion(ctx *echo.Context) error {
	service := service.GetRegionService()

	id := ctx.Param("regionID")
	if err := service.DeleteRegion(ctx, id); err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, nil)
}
