package repository

import (
	"apis/models"

	"github.com/labstack/echo/v5"
)

type AthleteRepository struct {
	QueryAthletes func(ctx *echo.Context) (*[]models.Athlete, error)
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

func queryAthletes(ctx *echo.Context) (*[]models.Athlete, error) {
	return &[]models.Athlete{}, nil
}

func putAthlete(ctx *echo.Context, athlete *models.Athlete) error {
	return nil
}

func getAthlete(ctx *echo.Context, athleteID string) (*models.Athlete, error) {
	return &models.Athlete{}, nil
}

func deleteAthlete(ctx *echo.Context, athleteID string) error {
	return nil
}
