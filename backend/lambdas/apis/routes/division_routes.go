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

func queryDivisions(ctx *echo.Context) error {
	service := service.GetDivisionService()

	page, err := service.QueryDivisions(ctx)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, page)
}

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

func deleteDivision(ctx *echo.Context) error {
	service := service.GetDivisionService()

	id := ctx.Param("divisionID")
	if err := service.DeleteDivision(ctx, id); err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, nil)
}
