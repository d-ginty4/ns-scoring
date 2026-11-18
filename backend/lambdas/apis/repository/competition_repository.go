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
	if db == nil {
		ConnectDB()
	}
	var competitions []models.Competition
	return &competitions, db.Preload("Region").Find(&competitions).Error
}

func putCompetition(ctx *echo.Context, competition *models.Competition) error {
	if db == nil {
		ConnectDB()
	}
	if err := db.Omit("Region").Save(competition).Error; err != nil {
		return err
	}
	return db.Preload("Region").First(competition, competition.ID).Error
}

func getCompetition(ctx *echo.Context, competitionID string) (*models.Competition, error) {
	if db == nil {
		ConnectDB()
	}

	var competition models.Competition
	result := db.Preload("Region").First(&competition, competitionID)
	if result.Error != nil {
		return nil, result.Error
	}

	return &competition, nil
}

func deleteCompetition(ctx *echo.Context, competitionID string) error {
	if db == nil {
		ConnectDB()
	}
	return (db.Unscoped().Delete(&models.Competition{}, competitionID)).Error
}
