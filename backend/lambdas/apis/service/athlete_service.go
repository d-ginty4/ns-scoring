package service

import (
	"fmt"

	"apis/exceptions"
	"apis/logger"
	"apis/models"
	"apis/repository"

	"github.com/labstack/echo/v5"
)

type AthleteService struct {
	QueryAthletes func(ctx *echo.Context) (*[]models.Athlete, error)
	CreateAthlete func(ctx *echo.Context, req *models.AthleteRequest, requestUser *models.User) (*models.Athlete, error)
	GetAthlete    func(ctx *echo.Context, id string) (*models.Athlete, error)
	UpdateAthlete func(ctx *echo.Context, athleteID string, req *models.AthleteRequest, existingAthlete *models.Athlete) (*models.Athlete, error)
	DeleteAthlete func(ctx *echo.Context, id string, existingAthlete *models.Athlete) error
}

var athleteService *AthleteService

func GetAthleteService() *AthleteService {
	if athleteService == nil {
		return newAthleteService()
	}
	return athleteService
}

func newAthleteService() *AthleteService {
	service := &AthleteService{
		QueryAthletes: queryAthletes,
		CreateAthlete: createAthlete,
		GetAthlete:    getAthlete,
		UpdateAthlete: updateAthlete,
		DeleteAthlete: deleteAthlete,
	}

	athleteService = service
	return service
}

func queryAthletes(ctx *echo.Context) (*[]models.Athlete, error) {
	athletes, err := repository.GetAthleteRepository().QueryAthletes(ctx)
	if err != nil {
		return nil, err
	}

	logger.Action(ctx, "query_athletes", "count", len(*athletes))
	return athletes, nil
}

func getAthlete(ctx *echo.Context, id string) (*models.Athlete, error) {
	athleteRepository := repository.GetAthleteRepository()

	athlete, err := athleteRepository.GetAthlete(ctx, id)
	if err != nil {
		return nil, err
	}

	if athlete == nil {
		return nil, exceptions.NewErrNotFound(fmt.Sprintf("Athlete with id %s not found", id))
	}

	logger.Action(ctx, "get_athlete", "athleteID", athlete.ID)
	return athlete, nil
}

func createAthlete(ctx *echo.Context, req *models.AthleteRequest, requestUser *models.User) (*models.Athlete, error) {
	repository := repository.GetAthleteRepository()

	athlete := models.Athlete{
		AthleteRequest: *req,
	}

	if err := repository.PutAthlete(ctx, &athlete); err != nil {
		return nil, err
	}

	logger.Action(ctx, "create_athlete", "athleteID", athlete.ID)
	return &athlete, nil
}

func updateAthlete(ctx *echo.Context, athleteID string, req *models.AthleteRequest, existingAthlete *models.Athlete) (*models.Athlete, error) {
	repository := repository.GetAthleteRepository()

	athlete := *existingAthlete
	athlete.AthleteRequest = *req

	if err := repository.PutAthlete(ctx, &athlete); err != nil {
		return nil, err
	}

	logger.Action(ctx, "update_athlete", "athleteID", athlete.ID)
	return &athlete, nil
}

func deleteAthlete(ctx *echo.Context, id string, existingAthlete *models.Athlete) error {
	athleteRepository := repository.GetAthleteRepository()

	if err := athleteRepository.DeleteAthlete(ctx, id); err != nil {
		return err
	}

	logger.Action(ctx, "delete_athlete", "athleteID", id)
	return nil
}
