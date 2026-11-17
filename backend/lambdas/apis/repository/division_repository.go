package repository

import (
	"apis/models"

	"github.com/labstack/echo/v5"
)

type DivisionRepository struct {
	QueryDivisions func(ctx *echo.Context) (*[]models.Division, error)
	PutDivision    func(ctx *echo.Context, division *models.Division) error
	GetDivision    func(ctx *echo.Context, divisionID string) (*models.Division, error)
	DeleteDivision func(ctx *echo.Context, divisionID string) error
}

var divisionRepository *DivisionRepository

func GetDivisionRepository() *DivisionRepository {
	if divisionRepository == nil {
		return newDivisionRepository()
	}

	return divisionRepository
}

func newDivisionRepository() *DivisionRepository {
	repository := &DivisionRepository{
		QueryDivisions: queryDivisions,
		PutDivision:    putDivision,
		GetDivision:    getDivision,
		DeleteDivision: deleteDivision,
	}

	divisionRepository = repository
	return repository
}

func queryDivisions(ctx *echo.Context) (*[]models.Division, error) {
	if db == nil {
		ConnectDB()
	}
	var divisions []models.Division
	return &divisions, (db.Find(&divisions)).Error
}

func getDivision(ctx *echo.Context, divisionID string) (*models.Division, error) {
	if db == nil {
		ConnectDB()
	}

	var division models.Division
	result := db.First(&division, divisionID)
	if result.Error != nil {
		return nil, result.Error
	}

	return &division, nil
}

func putDivision(ctx *echo.Context, division *models.Division) error {
	if db == nil {
		ConnectDB()
	}
	return db.Save(division).Error
}

func deleteDivision(ctx *echo.Context, divisionID string) error {
	if db == nil {
		ConnectDB()
	}
	return (db.Unscoped().Delete(&models.Division{}, divisionID)).Error
}
