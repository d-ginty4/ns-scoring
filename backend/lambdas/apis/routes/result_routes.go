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

	page, err := service.QueryResults(ctx)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, page)
}

func postResult(ctx *echo.Context) error {
	service := service.GetResultService()

	requestUser, ok := ctx.Get("requestUser").(*models.User)
	if !ok || requestUser == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get request user")
	}

	var p models.ResultRequest
	if err := utils.BindAndValidate(ctx, &p); err != nil {
		return err
	}

	result, err := service.CreateResult(ctx, &p, requestUser)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, result)
}

func getResult(ctx *echo.Context) error {
	resultService := service.GetResultService()

	id := ctx.Param("resultID")

	result, ok := ctx.Get("result").(*models.Result)
	if ok && result != nil {
		return responses.Success(ctx, http.StatusOK, result)
	}

	result, err := resultService.GetResult(ctx, id)
	if err != nil {
		if errors.Is(err, &exceptions.ErrNotFound{}) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return responses.Success(ctx, http.StatusOK, result)
}

func putResult(ctx *echo.Context) error {
	service := service.GetResultService()

	existing, ok := ctx.Get("result").(*models.Result)
	if !ok || existing == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get existing result")
	}

	var p models.ResultRequest
	if err := utils.BindAndValidate(ctx, &p); err != nil {
		return err
	}

	id := ctx.Param("resultID")
	result, err := service.UpdateResult(ctx, id, &p, existing)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, result)
}

func deleteResult(ctx *echo.Context) error {
	service := service.GetResultService()

	result, ok := ctx.Get("result").(*models.Result)
	if !ok || result == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get result")
	}

	id := ctx.Param("resultID")
	if err := service.DeleteResult(ctx, id, result); err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, nil)
}
