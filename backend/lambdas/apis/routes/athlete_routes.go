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

type AthleteRoutes struct {
	QueryAthletes func(c *echo.Context) error
	PostAthlete   func(c *echo.Context) error
	GetAthlete    func(c *echo.Context) error
	PutAthlete    func(c *echo.Context) error
	DeleteAthlete func(c *echo.Context) error
}

var athleteRoutes *AthleteRoutes

func getAthleteRoutes() *AthleteRoutes {
	if athleteRoutes == nil {
		return newAthleteRoutes()
	}
	return athleteRoutes
}

func newAthleteRoutes() *AthleteRoutes {
	routes := &AthleteRoutes{
		QueryAthletes: queryAthletes,
		PostAthlete:   postAthlete,
		GetAthlete:    getAthlete,
		PutAthlete:    putAthlete,
		DeleteAthlete: deleteAthlete,
	}

	athleteRoutes = routes
	return routes
}

func queryAthletes(ctx *echo.Context) error {
	service := service.GetAthleteService()

	page, err := service.QueryAthletes(ctx)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, page)
}

func postAthlete(ctx *echo.Context) error {
	service := service.GetAthleteService()

	requestUser, ok := ctx.Get("requestUser").(*models.User)
	if !ok || requestUser == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get request user")
	}

	var p models.AthleteRequest
	if err := utils.BindAndValidate(ctx, &p); err != nil {
		return err
	}

	athlete, err := service.CreateAthlete(ctx, &p, requestUser)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, athlete)
}

func getAthlete(ctx *echo.Context) error {
	athleteService := service.GetAthleteService()

	id := ctx.Param("athleteID")

	athlete, ok := ctx.Get("athlete").(*models.Athlete)
	if ok && athlete != nil {
		return responses.Success(ctx, http.StatusOK, athlete)
	}

	athlete, err := athleteService.GetAthlete(ctx, id)
	if err != nil {
		if errors.Is(err, &exceptions.ErrNotFound{}) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return responses.Success(ctx, http.StatusOK, athlete)
}

func putAthlete(ctx *echo.Context) error {
	service := service.GetAthleteService()

	existing, ok := ctx.Get("athlete").(*models.Athlete)
	if !ok || existing == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get existing athlete")
	}

	var p models.AthleteRequest
	if err := utils.BindAndValidate(ctx, &p); err != nil {
		return err
	}

	id := ctx.Param("athleteID")
	athlete, err := service.UpdateAthlete(ctx, id, &p, existing)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, athlete)
}

func deleteAthlete(ctx *echo.Context) error {
	service := service.GetAthleteService()

	athlete, ok := ctx.Get("athlete").(*models.Athlete)
	if !ok || athlete == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get athlete")
	}

	id := ctx.Param("athleteID")
	if err := service.DeleteAthlete(ctx, id, athlete); err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, nil)
}
