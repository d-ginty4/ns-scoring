package routes

import (
	"net/http"

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

// queryAthletes lists athletes in a competition.
//
//	@Summary		List athletes
//	@Tags			Athletes
//	@Produce		json
//	@Param			competitionID	path		integer	true	"Competition ID"
//	@Success		200				{object}	responses.Envelope[[]models.Athlete]
//	@Failure		404				{object}	responses.Envelope[any]
//	@Failure		500				{object}	responses.Envelope[any]
//	@Router			/v1/competitions/{competitionID}/athletes [get]
func queryAthletes(ctx *echo.Context) error {
	service := service.GetAthleteService()

	page, err := service.QueryAthletes(ctx, ctx.Param("competitionID"))
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, page)
}

// postAthlete creates a athlete.
//
//	@Summary		Create athlete
//	@Tags			Athletes
//	@Accept			json
//	@Produce		json
//	@Param			competitionID	path		integer					true	"Competition ID"
//	@Param			athlete			body		models.AthleteRequest	true	"Athlete"
//	@Success		200				{object}	responses.Envelope[models.Athlete]
//	@Failure		400				{object}	responses.Envelope[any]
//	@Failure		404				{object}	responses.Envelope[any]
//	@Failure		500				{object}	responses.Envelope[any]
//	@Router			/v1/competitions/{competitionID}/athletes [post]
func postAthlete(ctx *echo.Context) error {
	service := service.GetAthleteService()

	var p models.AthleteRequest
	if err := utils.BindAndValidate(ctx, &p); err != nil {
		return err
	}

	athlete, err := service.CreateAthlete(ctx, ctx.Param("competitionID"), &p)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, athlete)
}

// getAthlete returns an athlete.
//
//	@Summary		Get athlete
//	@Tags			Athletes
//	@Produce		json
//	@Param			competitionID	path		integer	true	"Competition ID"
//	@Param			athleteID		path		integer	true	"Athlete ID"
//	@Success		200				{object}	responses.Envelope[models.Athlete]
//	@Failure		404				{object}	responses.Envelope[any]
//	@Failure		500				{object}	responses.Envelope[any]
//	@Router			/v1/competitions/{competitionID}/athletes/{athleteID} [get]
func getAthlete(ctx *echo.Context) error {
	athlete, ok := ctx.Get("athlete").(*models.Athlete)
	if ok && athlete != nil {
		return responses.Success(ctx, http.StatusOK, athlete)
	}

	athleteService := service.GetAthleteService()
	athlete, err := athleteService.GetAthlete(ctx, ctx.Param("competitionID"), ctx.Param("athleteID"))
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, athlete)
}

// putAthlete creates or updates an athlete.
//
//	@Summary		Update athlete
//	@Tags			Athletes
//	@Accept			json
//	@Produce		json
//	@Param			competitionID	path		integer					true	"Competition ID"
//	@Param			athleteID		path		integer					true	"Athlete ID"
//	@Param			athlete			body		models.AthleteRequest	true	"Athlete"
//	@Success		200				{object}	responses.Envelope[models.Athlete]
//	@Failure		400				{object}	responses.Envelope[any]
//	@Failure		404				{object}	responses.Envelope[any]
//	@Failure		500				{object}	responses.Envelope[any]
//	@Router			/v1/competitions/{competitionID}/athletes/{athleteID} [put]
func putAthlete(ctx *echo.Context) error {
	service := service.GetAthleteService()

	var p models.AthleteRequest
	if err := utils.BindAndValidate(ctx, &p); err != nil {
		return err
	}

	existingAthlete, _ := ctx.Get("athlete").(*models.Athlete)
	athlete, err := service.UpdateAthlete(ctx, ctx.Param("competitionID"), ctx.Param("athleteID"), &p, existingAthlete)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, athlete)
}

// deleteAthlete deletes an athlete.
//
//	@Summary		Delete athlete
//	@Tags			Athletes
//	@Produce		json
//	@Param			competitionID	path		integer	true	"Competition ID"
//	@Param			athleteID		path		integer	true	"Athlete ID"
//	@Success		200				{object}	responses.Envelope[models.Athlete]
//	@Failure		404				{object}	responses.Envelope[any]
//	@Failure		500				{object}	responses.Envelope[any]
//	@Router			/v1/competitions/{competitionID}/athletes/{athleteID} [delete]
func deleteAthlete(ctx *echo.Context) error {
	service := service.GetAthleteService()

	if err := service.DeleteAthlete(ctx, ctx.Param("competitionID"), ctx.Param("athleteID")); err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, nil)
}
