package repository

import (
	"apis/models"

	"github.com/labstack/echo/v5"
)

type AthleteRepository struct {
	QueryAthletes func(ctx *echo.Context, competitionID string) (*[]models.Athlete, error)
	PutAthlete    func(ctx *echo.Context, athlete *models.Athlete) error
	GetAthlete    func(ctx *echo.Context, athleteID string) (*models.Athlete, error)
	DeleteAthlete func(ctx *echo.Context, athleteID string) error
}

var athleteRepository *AthleteRepository

func GetAthleteRepository() *AthleteRepository {
	if athleteRepository == nil {
		return newAthleteRepository()
	}

	return athleteRepository
}

func newAthleteRepository() *AthleteRepository {
	repository := &AthleteRepository{
		QueryAthletes: queryAthletes,
		PutAthlete:    putAthlete,
		GetAthlete:    getAthlete,
		DeleteAthlete: deleteAthlete,
	}

	athleteRepository = repository
	return repository
}

func queryAthletes(ctx *echo.Context, competitionID string) (*[]models.Athlete, error) {
	if db == nil {
		ConnectDB()
	}
	var athletes []models.Athlete
	return &athletes, db.Preload("Competition.Region").Preload("Division").Where("competition_id = ?", competitionID).Find(&athletes).Error
}

func putAthlete(ctx *echo.Context, athlete *models.Athlete) error {
	if db == nil {
		ConnectDB()
	}
	if err := db.Omit("Competition", "Division").Save(athlete).Error; err != nil {
		return err
	}
	return db.Preload("Competition.Region").Preload("Division").First(athlete, athlete.ID).Error
}

func getAthlete(ctx *echo.Context, athleteID string) (*models.Athlete, error) {
	if db == nil {
		ConnectDB()
	}

	var athlete models.Athlete
	result := db.Preload("Competition.Region").Preload("Division").First(&athlete, athleteID)
	if result.Error != nil {
		return nil, result.Error
	}

	return &athlete, nil
}

func deleteAthlete(ctx *echo.Context, athleteID string) error {
	if db == nil {
		ConnectDB()
	}
	return db.Unscoped().Delete(&models.Athlete{}, athleteID).Error
}
