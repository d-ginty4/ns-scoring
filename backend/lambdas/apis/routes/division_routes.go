package routes

import (
	"net/http"

	"apis/models"
	"apis/responses"
	"apis/service"
	"apis/utils"

	"github.com/labstack/echo/v5"
)

type DivisionRoutes struct {
	QueryDivisions func(c *echo.Context) error
	PostDivision   func(c *echo.Context) error
	GetDivision    func(c *echo.Context) error
	PutDivision    func(c *echo.Context) error
	DeleteDivision func(c *echo.Context) error
}

var divisionRoutes *DivisionRoutes

func getDivisionRoutes() *DivisionRoutes {
	if divisionRoutes == nil {
		return newDivisionRoutes()
	}
	return divisionRoutes
}

func newDivisionRoutes() *DivisionRoutes {
	routes := &DivisionRoutes{
		QueryDivisions: queryDivisions,
		PostDivision:   postDivision,
		GetDivision:    getDivision,
		PutDivision:    putDivision,
		DeleteDivision: deleteDivision,
	}

	divisionRoutes = routes
	return routes
}

// queryDivisions lists divisions.
//
//	@Summary		List divisions
//	@Tags			Divisions
//	@Produce		json
//	@Success		200	{object}	responses.Envelope[[]models.Division]
//	@Failure		500	{object}	responses.Envelope[any]
//	@Router			/v1/divisions [get]
func queryDivisions(ctx *echo.Context) error {
	service := service.GetDivisionService()

	page, err := service.QueryDivisions(ctx)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, page)
}

// postDivision creates a division.
//
//	@Summary		Create division
//	@Tags			Divisions
//	@Accept			json
//	@Produce		json
//	@Param			division	body		models.DivisionRequest	true	"Division"
//	@Success		200			{object}	responses.Envelope[models.Division]
//	@Failure		400			{object}	responses.Envelope[any]
//	@Failure		500			{object}	responses.Envelope[any]
//	@Router			/v1/divisions [post]
func postDivision(ctx *echo.Context) error {
	service := service.GetDivisionService()

	var p models.DivisionRequest
	if err := utils.BindAndValidate(ctx, &p); err != nil {
		return err
	}

	division, err := service.CreateDivision(ctx, &p)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, division)
}

// getDivision returns a division.
//
//	@Summary		Get division
//	@Tags			Divisions
//	@Produce		json
//	@Param			divisionID	path		integer	true	"Division ID"
//	@Success		200			{object}	responses.Envelope[models.Division]
//	@Failure		404			{object}	responses.Envelope[any]
//	@Failure		500			{object}	responses.Envelope[any]
//	@Router			/v1/divisions/{divisionID} [get]
func getDivision(ctx *echo.Context) error {
	division, ok := ctx.Get("division").(*models.Division)
	if ok && division != nil {
		return responses.Success(ctx, http.StatusOK, division)
	}

	divisionService := service.GetDivisionService()
	id := ctx.Param("divisionID")

	division, err := divisionService.GetDivision(ctx, id)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, division)
}

// putDivision creates or updates a division.
//
//	@Summary		Update division
//	@Tags			Divisions
//	@Accept			json
//	@Produce		json
//	@Param			divisionID	path		integer						true	"Division ID"
//	@Param			division	body		models.DivisionRequest		true	"Division"
//	@Success		200			{object}	responses.Envelope[models.Division]
//	@Failure		400			{object}	responses.Envelope[any]
//	@Failure		404			{object}	responses.Envelope[any]
//	@Failure		500			{object}	responses.Envelope[any]
//	@Router			/v1/divisions/{divisionID} [put]
func putDivision(ctx *echo.Context) error {
	service := service.GetDivisionService()

	var p models.DivisionRequest
	if err := utils.BindAndValidate(ctx, &p); err != nil {
		return err
	}

	id := ctx.Param("divisionID")
	existingDivision, _ := ctx.Get("division").(*models.Division)
	division, err := service.UpdateDivision(ctx, id, &p, existingDivision)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, division)
}

// deleteDivision deletes a division.
//
//	@Summary		Delete division
//	@Tags			Divisions
//	@Produce		json
//	@Param			divisionID	path		integer	true	"Division ID"
//	@Success		200			{object}	responses.Envelope[models.Division]
//	@Failure		404			{object}	responses.Envelope[any]
//	@Failure		500			{object}	responses.Envelope[any]
//	@Router			/v1/divisions/{divisionID} [delete]
func deleteDivision(ctx *echo.Context) error {
	service := service.GetDivisionService()

	id := ctx.Param("divisionID")
	if err := service.DeleteDivision(ctx, id); err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, nil)
}
