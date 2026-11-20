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

// queryResults lists an athlete's results.
//
//	@Summary		List results
//	@Tags			Results
//	@Produce		json
//	@Param			athleteID	path		integer	true	"Athlete ID"
//	@Success		200			{object}	responses.Envelope[[]models.Result]
//	@Failure		500			{object}	responses.Envelope[any]
//	@Router			/v1/results/athletes/{athleteID} [get]
func queryResults(ctx *echo.Context) error {
	service := service.GetResultService()

	page, err := service.QueryResults(ctx, ctx.Param("athleteID"))
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, page)
}

// postResult creates a result.
//
//	@Summary		Create result
//	@Tags			Results
//	@Accept			json
//	@Produce		json
//	@Param			athleteID	path		integer					true	"Athlete ID"
//	@Param			eventID		path		integer					true	"Event ID"
//	@Param			result		body		models.ResultRequest	true	"Result"
//	@Success		200			{object}	responses.Envelope[models.Result]
//	@Failure		400			{object}	responses.Envelope[any]
//	@Failure		404			{object}	responses.Envelope[any]
//	@Failure		500			{object}	responses.Envelope[any]
//	@Router			/v1/results/athletes/{athleteID}/events/{eventID} [post]
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

// getResult returns an athlete's result for an event.
//
//	@Summary		Get result
//	@Tags			Results
//	@Produce		json
//	@Param			athleteID	path		integer	true	"Athlete ID"
//	@Param			eventID		path		integer	true	"Event ID"
//	@Success		200			{object}	responses.Envelope[models.Result]
//	@Failure		404			{object}	responses.Envelope[any]
//	@Failure		500			{object}	responses.Envelope[any]
//	@Router			/v1/results/athletes/{athleteID}/events/{eventID} [get]
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

// putResult creates or updates a result.
//
//	@Summary		Update result
//	@Tags			Results
//	@Accept			json
//	@Produce		json
//	@Param			athleteID	path		integer					true	"Athlete ID"
//	@Param			eventID		path		integer					true	"Event ID"
//	@Param			result		body		models.ResultRequest	true	"Result"
//	@Success		200			{object}	responses.Envelope[models.Result]
//	@Failure		400			{object}	responses.Envelope[any]
//	@Failure		404			{object}	responses.Envelope[any]
//	@Failure		500			{object}	responses.Envelope[any]
//	@Router			/v1/results/athletes/{athleteID}/events/{eventID} [put]
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

// deleteResult deletes an athlete's result for an event.
//
//	@Summary		Delete result
//	@Tags			Results
//	@Produce		json
//	@Param			athleteID	path		integer	true	"Athlete ID"
//	@Param			eventID		path		integer	true	"Event ID"
//	@Success		200			{object}	responses.Envelope[models.Result]
//	@Failure		404			{object}	responses.Envelope[any]
//	@Failure		500			{object}	responses.Envelope[any]
//	@Router			/v1/results/athletes/{athleteID}/events/{eventID} [delete]
func deleteResult(ctx *echo.Context) error {
	service := service.GetResultService()

	if err := service.DeleteResult(ctx, ctx.Param("athleteID"), ctx.Param("eventID")); err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, nil)
}
