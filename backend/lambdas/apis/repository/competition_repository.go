package repository

import (
	"apis/models"

	"github.com/labstack/echo/v5"
)

type CompetitionRepository struct {
	QueryCompetitions func(ctx *echo.Context) (*[]models.Competition, error)
	PutCompetition    func(ctx *echo.Context, competition *models.Competition) error
	GetCompetition    func(ctx *echo.Context, competitionID string) (*models.Competition, error)
	DeleteCompetition func(ctx *echo.Context, competitionID string) error
}

var competitionRepository *CompetitionRepository

func GetCompetitionRepository() *CompetitionRepository {
	if competitionRepository == nil {
		return newCompetitionRepository()
	}

	return competitionRepository
}

func newCompetitionRepository() *CompetitionRepository {
	repository := &CompetitionRepository{
		QueryCompetitions: queryCompetitions,
		PutCompetition:    putCompetition,
		GetCompetition:    getCompetition,
		DeleteCompetition: deleteCompetition,
	}

	competitionRepository = repository
	return repository
}

func queryCompetitions(ctx *echo.Context) (*[]models.Competition, error) {
	return &[]models.Competition{}, nil
}

func putCompetition(ctx *echo.Context, competition *models.Competition) error {
	return nil
}

func getCompetition(ctx *echo.Context, competitionID string) (*models.Competition, error) {
	return &models.Competition{}, nil
}

func deleteCompetition(ctx *echo.Context, competitionID string) error {
	return nil
}
