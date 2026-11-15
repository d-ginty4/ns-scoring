package repository

import (
	"apis/models"

	"github.com/labstack/echo/v5"
)

type ResultRepository struct {
	QueryResults func(ctx *echo.Context) (*[]models.Result, error)
	PutResult    func(ctx *echo.Context, result *models.Result) error
	GetResult    func(ctx *echo.Context, resultID string) (*models.Result, error)
	DeleteResult func(ctx *echo.Context, resultID string) error
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
		QueryResults: queryResults,
		PutResult:    putResult,
		GetResult:    getResult,
		DeleteResult: deleteResult,
	}

	resultRepository = repository
	return repository
}

func queryResults(ctx *echo.Context) (*[]models.Result, error) {
	return &[]models.Result{}, nil
}

func putResult(ctx *echo.Context, result *models.Result) error {
	return nil
}

func getResult(ctx *echo.Context, resultID string) (*models.Result, error) {
	return &models.Result{}, nil
}

func deleteResult(ctx *echo.Context, resultID string) error {
	return nil
}
