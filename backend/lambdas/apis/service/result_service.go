package service

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

	"apis/exceptions"
	"apis/logger"
	"apis/models"
	"apis/repository"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type ResultService struct {
	QueryResults         func(ctx *echo.Context, athleteID string) (*[]models.Result, error)
	CreateResult         func(ctx *echo.Context, athleteID string, eventID string, req *models.ResultRequest) (*models.Result, error)
	CreateAthleteResults func(ctx *echo.Context, athleteID string, req *models.AthleteResultsRequest) (*[]models.Result, error)
	GetResult            func(ctx *echo.Context, athleteID string, eventID string) (*models.Result, error)
	UpdateResult         func(ctx *echo.Context, athleteID string, eventID string, req *models.ResultRequest, existingResult *models.Result) (*models.Result, error)
	DeleteResult         func(ctx *echo.Context, athleteID string, eventID string) error
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
		QueryResults:         queryResults,
		CreateResult:         createResult,
		CreateAthleteResults: createAthleteResults,
		GetResult:            getResult,
		UpdateResult:         updateResult,
		DeleteResult:         deleteResult,
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

func createAthleteResults(ctx *echo.Context, athleteID string, req *models.AthleteResultsRequest) (*[]models.Result, error) {
	resultRepository := repository.GetResultRepository()

	athlete, ok := ctx.Get("athlete").(*models.Athlete)
	if !ok || athlete == nil {
		panic("athlete no exist")
	}

	events, err := repository.GetEventRepository().QueryEvents(ctx)
	if err != nil {
		return nil, err
	}
	eventsByID := make(map[uint]models.Event, len(*events))
	for _, event := range *events {
		eventsByID[event.ID] = event
	}

	existingResult, err := resultRepository.QueryResults(ctx, athleteID)
	if err != nil {
		return nil, err
	}
	existingEvents := make(map[uint]struct{}, len(*existingResult))
	for _, result := range *existingResult {
		existingEvents[result.EventId] = struct{}{}
	}

	var requestErrors []string
	var seen []uint
	results := make([]models.Result, 0, len(req.Results))
	for i, resultRequest := range req.Results {
		if slices.Contains(seen, resultRequest.EventID) {
			requestErrors = append(requestErrors, fmt.Sprintf("results[%d]: event %d is repeated", i, resultRequest.EventID))
			continue
		}
		seen = append(seen, resultRequest.EventID)

		event, ok := eventsByID[resultRequest.EventID]
		if !ok {
			requestErrors = append(requestErrors, fmt.Sprintf("results[%d]: event %d does not exist", i, resultRequest.EventID))
			continue
		}
		if _, ok := existingEvents[resultRequest.EventID]; ok {
			requestErrors = append(requestErrors, fmt.Sprintf("results[%d]: result for athlete %s and event %d already exists", i, athleteID, resultRequest.EventID))
			continue
		}
		if err := checkResultWeight(athlete, &event, resultRequest.Weight); err != nil {
			requestErrors = append(requestErrors, fmt.Sprintf("results[%d]: %s", i, err.Error()))
			continue
		}

		// If there's errors there is not point constructing any more results, they won't be used
		if len(requestErrors) > 0 {
			continue
		}

		results = append(results, models.Result{
			ResultRequest: resultRequest.ResultRequest,
			AthleteId:     athlete.ID,
			EventId:       event.ID,
		})
	}

	if len(requestErrors) > 0 {
		return nil, &exceptions.ErrValidation{Errors: requestErrors}
	}

	if err := resultRepository.CreateResults(ctx, &results); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, exceptions.NewErrBadEntity(fmt.Sprintf("result for athlete %s already exists", athleteID))
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, exceptions.NewErrBadEntity(fmt.Sprintf("athlete with id %s does not exist", athleteID))
		}
		return nil, err
	}

	logger.Action(ctx, "create_athlete_results", "athleteID", athleteID, "count", len(results))
	return &results, nil
}

func createResult(ctx *echo.Context, athleteID string, eventID string, req *models.ResultRequest) (*models.Result, error) {
	resultRepository := repository.GetResultRepository()

	athlete, ok := ctx.Get("athlete").(*models.Athlete)
	if !ok || athlete == nil {
		panic("athlete no exist")
	}
	event, ok := ctx.Get("event").(*models.Event)
	if !ok || event == nil {
		panic("event no exist")
	}

	if err := checkResultWeight(athlete, event, req.Weight); err != nil {
		return nil, err
	}

	if _, err := resultRepository.GetResult(ctx, athleteID, eventID); err == nil {
		return nil, exceptions.NewErrBadEntity(fmt.Sprintf("result for athlete %s and event %s already exists", athleteID, eventID))
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	result := models.Result{
		ResultRequest: *req,
		AthleteId:     athlete.ID,
		EventId:       event.ID,
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

	athlete, ok := ctx.Get("athlete").(*models.Athlete)
	if !ok || athlete == nil {
		panic("athlete no exist")
	}
	event, ok := ctx.Get("event").(*models.Event)
	if !ok || event == nil {
		panic("event no exist")
	}

	if err := checkResultWeight(athlete, event, req.Weight); err != nil {
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
	result.AthleteId = athlete.ID
	result.EventId = event.ID

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

// Check that the weight in the result is one of the options in this event
func checkResultWeight(athlete *models.Athlete, event *models.Event, weight int) error {
	var allowedWeights []int
	if athlete.Division.Gender == "male" {
		allowedWeights = event.WeightOptionsMen
	} else if athlete.Division.Gender == "female" {
		allowedWeights = event.WeightOptionsWomen
	} else {
		return exceptions.NewErrBadEntity("athlete division has no gender")
	}
	if !slices.Contains(allowedWeights, weight) {
		return exceptions.NewErrBadEntity(fmt.Sprintf(
			"weight %d is not one of the allowed %s options for this event %d: %s - %v",
			weight, athlete.Division.Gender, event.EventNumber, event.Name, allowedWeights,
		))
	}
	return nil
}
