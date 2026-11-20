package routes

import (
	"net/http"

	"apis/models"
	"apis/responses"
	"apis/service"
	"apis/utils"

	"github.com/labstack/echo/v5"
)

type CompetitionRoutes struct {
	QueryCompetitions func(c *echo.Context) error
	PostCompetition   func(c *echo.Context) error
	GetCompetition    func(c *echo.Context) error
	PutCompetition    func(c *echo.Context) error
	DeleteCompetition func(c *echo.Context) error
}

var competitionRoutes *CompetitionRoutes

func getCompetitionRoutes() *CompetitionRoutes {
	if competitionRoutes == nil {
		return newCompetitionRoutes()
	}
	return competitionRoutes
}

func newCompetitionRoutes() *CompetitionRoutes {
	routes := &CompetitionRoutes{
		QueryCompetitions: queryCompetitions,
		PostCompetition:   postCompetition,
		GetCompetition:    getCompetition,
		PutCompetition:    putCompetition,
		DeleteCompetition: deleteCompetition,
	}

	competitionRoutes = routes
	return routes
}

// queryCompetitions lists competitions.
//
//	@Summary		List competitions
//	@Tags			Competitions
//	@Produce		json
//	@Success		200	{object}	responses.Envelope[[]models.Competition]
//	@Failure		500	{object}	responses.Envelope[any]
//	@Router			/v1/competitions [get]
func queryCompetitions(ctx *echo.Context) error {
	service := service.GetCompetitionService()

	page, err := service.QueryCompetitions(ctx)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, page)
}

// postCompetition creates a competition.
//
//	@Summary		Create competition
//	@Tags			Competitions
//	@Accept			json
//	@Produce		json
//	@Param			competition	body		models.CompetitionRequest	true	"Competition"
//	@Success		200			{object}	responses.Envelope[models.Competition]
//	@Failure		400			{object}	responses.Envelope[any]
//	@Failure		500			{object}	responses.Envelope[any]
//	@Router			/v1/competitions [post]
func postCompetition(ctx *echo.Context) error {
	service := service.GetCompetitionService()

	var p models.CompetitionRequest
	if err := utils.BindAndValidate(ctx, &p); err != nil {
		return err
	}

	competition, err := service.CreateCompetition(ctx, &p)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, competition)
}

// getCompetition returns a competition.
//
//	@Summary		Get competition
//	@Tags			Competitions
//	@Produce		json
//	@Param			competitionID	path		integer	true	"Competition ID"
//	@Success		200				{object}	responses.Envelope[models.Competition]
//	@Failure		404				{object}	responses.Envelope[any]
//	@Failure		500				{object}	responses.Envelope[any]
//	@Router			/v1/competitions/{competitionID} [get]
func getCompetition(ctx *echo.Context) error {
	competition, ok := ctx.Get("competition").(*models.Competition)
	if ok && competition != nil {
		return responses.Success(ctx, http.StatusOK, competition)
	}

	competitionService := service.GetCompetitionService()
	id := ctx.Param("competitionID")

	competition, err := competitionService.GetCompetition(ctx, id)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, competition)
}

// putCompetition creates or updates a competition.
//
//	@Summary		Update competition
//	@Tags			Competitions
//	@Accept			json
//	@Produce		json
//	@Param			competitionID	path		integer							true	"Competition ID"
//	@Param			competition		body		models.CompetitionRequest		true	"Competition"
//	@Success		200				{object}	responses.Envelope[models.Competition]
//	@Failure		400				{object}	responses.Envelope[any]
//	@Failure		404				{object}	responses.Envelope[any]
//	@Failure		500				{object}	responses.Envelope[any]
//	@Router			/v1/competitions/{competitionID} [put]
func putCompetition(ctx *echo.Context) error {
	service := service.GetCompetitionService()

	var p models.CompetitionRequest
	if err := utils.BindAndValidate(ctx, &p); err != nil {
		return err
	}

	id := ctx.Param("competitionID")
	existingCompetition, _ := ctx.Get("competition").(*models.Competition)
	competition, err := service.UpdateCompetition(ctx, id, &p, existingCompetition)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, competition)
}

// deleteCompetition deletes a competition.
//
//	@Summary		Delete competition
//	@Tags			Competitions
//	@Produce		json
//	@Param			competitionID	path		integer	true	"Competition ID"
//	@Success		200				{object}	responses.Envelope[models.Competition]
//	@Failure		404				{object}	responses.Envelope[any]
//	@Failure		500				{object}	responses.Envelope[any]
//	@Router			/v1/competitions/{competitionID} [delete]
func deleteCompetition(ctx *echo.Context) error {
	service := service.GetCompetitionService()

	id := ctx.Param("competitionID")
	if err := service.DeleteCompetition(ctx, id); err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, nil)
}
