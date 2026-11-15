package routes

import (
	"errors"
	"net/http"

	"apis/exceptions"
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

	requestUser, ok := ctx.Get("requestUser").(*models.User)
	if !ok || requestUser == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get request user")
	}

	var p models.DivisionRequest
	if err := utils.BindAndValidate(ctx, &p); err != nil {
		return err
	}

	division, err := service.CreateDivision(ctx, &p, requestUser)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, division)
}

func getDivision(ctx *echo.Context) error {
	divisionService := service.GetDivisionService()

	id := ctx.Param("divisionID")

	division, ok := ctx.Get("division").(*models.Division)
	if ok && division != nil {
		return responses.Success(ctx, http.StatusOK, division)
	}

	division, err := divisionService.GetDivision(ctx, id)
	if err != nil {
		if errors.Is(err, &exceptions.ErrNotFound{}) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return responses.Success(ctx, http.StatusOK, division)
}

func putDivision(ctx *echo.Context) error {
	service := service.GetDivisionService()

	existing, ok := ctx.Get("division").(*models.Division)
	if !ok || existing == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get existing division")
	}

	var p models.DivisionRequest
	if err := utils.BindAndValidate(ctx, &p); err != nil {
		return err
	}

	id := ctx.Param("divisionID")
	division, err := service.UpdateDivision(ctx, id, &p, existing)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, division)
}

func deleteDivision(ctx *echo.Context) error {
	service := service.GetDivisionService()

	division, ok := ctx.Get("division").(*models.Division)
	if !ok || division == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get division")
	}

	id := ctx.Param("divisionID")
	if err := service.DeleteDivision(ctx, id, division); err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, nil)
}
