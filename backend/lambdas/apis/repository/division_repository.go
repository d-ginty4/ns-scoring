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
	return &[]models.Division{}, nil
}

func putDivision(ctx *echo.Context, division *models.Division) error {
	return nil
}

func getDivision(ctx *echo.Context, divisionID string) (*models.Division, error) {
	return &models.Division{}, nil
}

func deleteDivision(ctx *echo.Context, divisionID string) error {
	return nil
}
