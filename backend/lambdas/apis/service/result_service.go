package service

import (
	"fmt"

	"apis/exceptions"
	"apis/logger"
	"apis/models"
	"apis/repository"

	"github.com/labstack/echo/v5"
)

type ResultService struct {
	QueryResults func(ctx *echo.Context) (*[]models.Result, error)
	CreateResult func(ctx *echo.Context, req *models.ResultRequest, requestUser *models.User) (*models.Result, error)
	GetResult    func(ctx *echo.Context, id string) (*models.Result, error)
	UpdateResult func(ctx *echo.Context, resultID string, req *models.ResultRequest, existingResult *models.Result) (*models.Result, error)
	DeleteResult func(ctx *echo.Context, id string, existingResult *models.Result) error
}

var resultService *ResultService

func GetResultService() *ResultService {
	if resultService == nil {
		return newResultService()
	}
	return resultService
}

func newResultService() *ResultService {
	service := &ResultService{
		QueryResults: queryResults,
		CreateResult: createResult,
		GetResult:    getResult,
		UpdateResult: updateResult,
		DeleteResult: deleteResult,
	}

	resultService = service
	return service
}

func queryResults(ctx *echo.Context) (*[]models.Result, error) {
	results, err := repository.GetResultRepository().QueryResults(ctx)
	if err != nil {
		return nil, err
	}

	logger.Action(ctx, "query_results", "count", len(*results))
	return results, nil
}

func getResult(ctx *echo.Context, id string) (*models.Result, error) {
	resultRepository := repository.GetResultRepository()

	result, err := resultRepository.GetResult(ctx, id)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, exceptions.NewErrNotFound(fmt.Sprintf("Result with id %s not found", id))
	}

	logger.Action(ctx, "get_result", "resultID", result.ID)
	return result, nil
}

func createResult(ctx *echo.Context, req *models.ResultRequest, requestUser *models.User) (*models.Result, error) {
	repository := repository.GetResultRepository()

	result := models.Result{
		ResultRequest: *req,
	}

	if err := repository.PutResult(ctx, &result); err != nil {
		return nil, err
	}

	logger.Action(ctx, "create_result", "resultID", result.ID)
	return &result, nil
}

func updateResult(ctx *echo.Context, resultID string, req *models.ResultRequest, existingResult *models.Result) (*models.Result, error) {
	repository := repository.GetResultRepository()

	result := *existingResult
	result.ResultRequest = *req

	if err := repository.PutResult(ctx, &result); err != nil {
		return nil, err
	}

	logger.Action(ctx, "update_result", "resultID", result.ID)
	return &result, nil
}

func deleteResult(ctx *echo.Context, id string, existingResult *models.Result) error {
	resultRepository := repository.GetResultRepository()

	if err := resultRepository.DeleteResult(ctx, id); err != nil {
		return err
	}

	logger.Action(ctx, "delete_result", "resultID", id)
	return nil
}
