package repository

import (
	"apis/models"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type ResultRepository struct {
	QueryResults  func(ctx *echo.Context, athleteID string) (*[]models.Result, error)
	PutResult     func(ctx *echo.Context, result *models.Result) error
	CreateResults func(ctx *echo.Context, results *[]models.Result) error
	GetResult     func(ctx *echo.Context, athleteID string, eventID string) (*models.Result, error)
	DeleteResult  func(ctx *echo.Context, athleteID string, eventID string) error
}

var resultRepository *ResultRepository

func GetResultRepository() *ResultRepository {
	if resultRepository == nil {
		return newResultRepository()
	}

	return resultRepository
}

func newResultRepository() *ResultRepository {
	repository := &ResultRepository{
		QueryResults:  queryResults,
		PutResult:     putResult,
		CreateResults: createResults,
		GetResult:     getResult,
		DeleteResult:  deleteResult,
	}

	resultRepository = repository
	return repository
}

func queryResults(ctx *echo.Context, athleteID string) (*[]models.Result, error) {
	if db == nil {
		ConnectDB()
	}
	var results []models.Result
	return &results, db.Preload("Event").
		Preload("Athlete.Competition.Region").
		Preload("Athlete.Division").
		Where("athlete_id = ?", athleteID).
		Find(&results).Error
}

func getResult(ctx *echo.Context, athleteID string, eventID string) (*models.Result, error) {
	if db == nil {
		ConnectDB()
	}

	var result models.Result
	query := db.Preload("Event").
		Preload("Athlete.Competition.Region").
		Preload("Athlete.Division").
		Where("athlete_id = ? AND event_id = ?", athleteID, eventID).
		First(&result)
	if query.Error != nil {
		return nil, query.Error
	}

	return &result, nil
}

func createResults(ctx *echo.Context, results *[]models.Result) error {
	if db == nil {
		ConnectDB()
	}
	if len(*results) == 0 {
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Event", "Athlete").Create(results).Error; err != nil {
			return err
		}

		ids := make([]uint, len(*results))
		for i, result := range *results {
			ids[i] = result.ID
		}

		var loaded []models.Result
		if err := tx.Preload("Event").
			Preload("Athlete.Competition.Region").
			Preload("Athlete.Division").
			Find(&loaded, ids).Error; err != nil {
			return err
		}

		*results = loaded
		return nil
	})
}

func putResult(ctx *echo.Context, result *models.Result) error {
	if db == nil {
		ConnectDB()
	}
	if err := db.Omit("Event", "Athlete").Save(result).Error; err != nil {
		return err
	}
	return db.Preload("Event").
		Preload("Athlete.Competition.Region").
		Preload("Athlete.Division").
		First(result, result.ID).Error
}

func deleteResult(ctx *echo.Context, athleteID string, eventID string) error {
	if db == nil {
		ConnectDB()
	}
	return db.Unscoped().
		Where("athlete_id = ? AND event_id = ?", athleteID, eventID).
		Delete(&models.Result{}).Error
}
