package routes

import (
	"net/http"

	"apis/models"
	"apis/responses"
	"apis/service"
	"apis/utils"

	"github.com/labstack/echo/v5"
)

type ResultRoutes struct {
	QueryResults func(c *echo.Context) error
	PostResult   func(c *echo.Context) error
	GetResult    func(c *echo.Context) error
	PutResult    func(c *echo.Context) error
	DeleteResult func(c *echo.Context) error
}

var resultRoutes *ResultRoutes

func getResultRoutes() *ResultRoutes {
	if resultRoutes == nil {
		return newResultRoutes()
	}
	return resultRoutes
}

func newResultRoutes() *ResultRoutes {
	routes := &ResultRoutes{
		QueryResults: queryResults,
		PostResult:   postResult,
		GetResult:    getResult,
		PutResult:    putResult,
		DeleteResult: deleteResult,
	}

	resultRoutes = routes
	return routes
}

func queryResults(ctx *echo.Context) error {
	service := service.GetResultService()

	page, err := service.QueryResults(ctx, ctx.Param("athleteID"))
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, page)
}

func postResult(ctx *echo.Context) error {
	service := service.GetResultService()

	var p models.ResultRequest
	if err := utils.BindAndValidate(ctx, &p); err != nil {
		return err
	}

	result, err := service.CreateResult(ctx, ctx.Param("athleteID"), ctx.Param("eventID"), &p)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, result)
}

func getResult(ctx *echo.Context) error {
	result, ok := ctx.Get("result").(*models.Result)
	if ok && result != nil {
		return responses.Success(ctx, http.StatusOK, result)
	}

	resultService := service.GetResultService()
	result, err := resultService.GetResult(ctx, ctx.Param("athleteID"), ctx.Param("eventID"))
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, result)
}

func putResult(ctx *echo.Context) error {
	service := service.GetResultService()

	var p models.ResultRequest
	if err := utils.BindAndValidate(ctx, &p); err != nil {
		return err
	}

	existingResult, _ := ctx.Get("result").(*models.Result)
	result, err := service.UpdateResult(ctx, ctx.Param("athleteID"), ctx.Param("eventID"), &p, existingResult)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, result)
}

func deleteResult(ctx *echo.Context) error {
	service := service.GetResultService()

	if err := service.DeleteResult(ctx, ctx.Param("athleteID"), ctx.Param("eventID")); err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, nil)
}
