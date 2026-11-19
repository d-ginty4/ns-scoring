package service

import (
	"errors"
	"fmt"
	"strconv"

	"apis/exceptions"
	"apis/logger"
	"apis/models"
	"apis/repository"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type ResultService struct {
	QueryResults func(ctx *echo.Context, athleteID string) (*[]models.Result, error)
	CreateResult func(ctx *echo.Context, athleteID string, eventID string, req *models.ResultRequest) (*models.Result, error)
	GetResult    func(ctx *echo.Context, athleteID string, eventID string) (*models.Result, error)
	UpdateResult func(ctx *echo.Context, athleteID string, eventID string, req *models.ResultRequest, existingResult *models.Result) (*models.Result, error)
	DeleteResult func(ctx *echo.Context, athleteID string, eventID string) error
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

func queryResults(ctx *echo.Context, athleteID string) (*[]models.Result, error) {
	results, err := repository.GetResultRepository().QueryResults(ctx, athleteID)
	if err != nil {
		return nil, err
	}

	logger.Action(ctx, "query_results", "athleteID", athleteID, "count", len(*results))
	return results, nil
}

func getResult(ctx *echo.Context, athleteID string, eventID string) (*models.Result, error) {
	resultRepository := repository.GetResultRepository()

	result, err := resultRepository.GetResult(ctx, athleteID, eventID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewErrNotFound(fmt.Sprintf("result for athlete %s and event %s not found", athleteID, eventID))
		}
		return nil, err
	}

	logger.Action(ctx, "get_result", "athleteID", athleteID, "eventID", eventID, "resultID", result.ID)
	return result, nil
}

func createResult(ctx *echo.Context, athleteID string, eventID string, req *models.ResultRequest) (*models.Result, error) {
	resultRepository := repository.GetResultRepository()

	athlete, err := idFromPath(athleteID, "athlete")
	if err != nil {
		return nil, err
	}
	event, err := idFromPath(eventID, "event")
	if err != nil {
		return nil, err
	}

	if _, err := resultRepository.GetResult(ctx, athleteID, eventID); err == nil {
		return nil, exceptions.NewErrBadEntity(fmt.Sprintf("result for athlete %s and event %s already exists", athleteID, eventID))
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	result := models.Result{
		ResultRequest: *req,
		AthleteId:     athlete,
		EventId:       event,
	}

	if err := resultRepository.PutResult(ctx, &result); err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, exceptions.NewErrBadEntity(fmt.Sprintf("athlete with id %s does not exist", athleteID))
		}
		return nil, err
	}

	logger.Action(ctx, "create_result", "athleteID", athleteID, "eventID", eventID, "resultID", result.ID)
	return &result, nil
}

func updateResult(ctx *echo.Context, athleteID string, eventID string, req *models.ResultRequest, result *models.Result) (*models.Result, error) {
	resultRepository := repository.GetResultRepository()

	athlete, err := idFromPath(athleteID, "athlete")
	if err != nil {
		return nil, err
	}
	event, err := idFromPath(eventID, "event")
	if err != nil {
		return nil, err
	}

	if result == nil {
		existing, err := resultRepository.GetResult(ctx, athleteID, eventID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, exceptions.NewErrNotFound(fmt.Sprintf("result for athlete %s and event %s not found", athleteID, eventID))
			}
			return nil, err
		}
		result = existing
	}
	result.ResultRequest = *req
	result.AthleteId = athlete
	result.EventId = event

	if err := resultRepository.PutResult(ctx, result); err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, exceptions.NewErrBadEntity(fmt.Sprintf("athlete with id %s does not exist", athleteID))
		}
		return nil, err
	}

	logger.Action(ctx, "update_result", "athleteID", athleteID, "eventID", eventID, "resultID", result.ID)
	return result, nil
}

func deleteResult(ctx *echo.Context, athleteID string, eventID string) error {
	resultRepository := repository.GetResultRepository()

	if err := resultRepository.DeleteResult(ctx, athleteID, eventID); err != nil {
		return err
	}

	logger.Action(ctx, "delete_result", "athleteID", athleteID, "eventID", eventID)
	return nil
}

func idFromPath(id string, name string) (uint, error) {
	parsed, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, exceptions.NewErrNotFound(fmt.Sprintf("%s with id %s not found", name, id))
	}
	return uint(parsed), nil
}
