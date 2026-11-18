package service

import (
	"errors"
	"fmt"

	"apis/exceptions"
	"apis/logger"
	"apis/models"
	"apis/repository"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type CompetitionService struct {
	QueryCompetitions func(ctx *echo.Context) (*[]models.Competition, error)
	CreateCompetition func(ctx *echo.Context, req *models.CompetitionRequest) (*models.Competition, error)
	GetCompetition    func(ctx *echo.Context, id string) (*models.Competition, error)
	UpdateCompetition func(ctx *echo.Context, competitionID string, req *models.CompetitionRequest, existingCompetition *models.Competition) (*models.Competition, error)
	DeleteCompetition func(ctx *echo.Context, id string) error
}

var competitionService *CompetitionService

func GetCompetitionService() *CompetitionService {
	if competitionService == nil {
		return newCompetitionService()
	}
	return competitionService
}

func newCompetitionService() *CompetitionService {
	service := &CompetitionService{
		QueryCompetitions: queryCompetitions,
		CreateCompetition: createCompetition,
		GetCompetition:    getCompetition,
		UpdateCompetition: updateCompetition,
		DeleteCompetition: deleteCompetition,
	}

	competitionService = service
	return service
}

func queryCompetitions(ctx *echo.Context) (*[]models.Competition, error) {
	competitions, err := repository.GetCompetitionRepository().QueryCompetitions(ctx)
	if err != nil {
		return nil, err
	}

	logger.Action(ctx, "query_competitions", "count", len(*competitions))
	return competitions, nil
}

func getCompetition(ctx *echo.Context, id string) (*models.Competition, error) {
	competitionRepository := repository.GetCompetitionRepository()

	competition, err := competitionRepository.GetCompetition(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewErrNotFound(fmt.Sprintf("competition with id %s not found", id))
		}
		return nil, err
	}

	logger.Action(ctx, "get_competition", "competitionID", competition.ID)
	return competition, nil
}

func createCompetition(ctx *echo.Context, req *models.CompetitionRequest) (*models.Competition, error) {
	repository := repository.GetCompetitionRepository()

	competition := models.Competition{
		CompetitionRequest: *req,
	}

	if err := repository.PutCompetition(ctx, &competition); err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, exceptions.NewErrBadEntity(fmt.Sprintf("region with id %d does not exist", competition.RegionID))
		}
		return nil, err
	}

	logger.Action(ctx, "create_competition", "competitionID", competition.ID)
	return &competition, nil
}

func updateCompetition(ctx *echo.Context, competitionID string, req *models.CompetitionRequest, competition *models.Competition) (*models.Competition, error) {
	repository := repository.GetCompetitionRepository()

	if competition == nil {
		c, err := repository.GetCompetition(ctx, competitionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, exceptions.NewErrNotFound(fmt.Sprintf("competition with id %s not found", competitionID))
			}
			return nil, err
		}
		competition = c
	}
	competition.CompetitionRequest = *req

	if err := repository.PutCompetition(ctx, competition); err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, exceptions.NewErrBadEntity(fmt.Sprintf("region with id %d does not exist", competition.RegionID))
		}
		return nil, err
	}

	logger.Action(ctx, "update_competition", "competitionID", competition.ID)
	return competition, nil
}

func deleteCompetition(ctx *echo.Context, id string) error {
	competitionRepository := repository.GetCompetitionRepository()

	if err := competitionRepository.DeleteCompetition(ctx, id); err != nil {
		return err
	}

	logger.Action(ctx, "delete_competition", "competitionID", id)
	return nil
}
