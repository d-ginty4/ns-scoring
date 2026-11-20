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

type AthleteService struct {
	QueryAthletes func(ctx *echo.Context, competitionID string) (*[]models.Athlete, error)
	CreateAthlete func(ctx *echo.Context, competitionID string, req *models.AthleteRequest) (*models.Athlete, error)
	GetAthlete    func(ctx *echo.Context, competitionID string, id string) (*models.Athlete, error)
	UpdateAthlete func(ctx *echo.Context, competitionID string, athleteID string, req *models.AthleteRequest, existingAthlete *models.Athlete) (*models.Athlete, error)
	DeleteAthlete func(ctx *echo.Context, competitionID string, id string) error
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

func queryAthletes(ctx *echo.Context, competitionID string) (*[]models.Athlete, error) {
	athletes, err := repository.GetAthleteRepository().QueryAthletes(ctx, competitionID)
	if err != nil {
		return nil, err
	}

	logger.Action(ctx, "query_athletes", "competitionID", competitionID, "count", len(*athletes))
	return athletes, nil
}

func getAthlete(ctx *echo.Context, competitionID string, id string) (*models.Athlete, error) {
	athleteRepository := repository.GetAthleteRepository()

	athlete, err := athleteRepository.GetAthlete(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewErrNotFound(fmt.Sprintf("athlete with id %s not found", id))
		}
		return nil, err
	}

	logger.Action(ctx, "get_athlete", "competitionID", competitionID, "athleteID", athlete.ID)
	return athlete, nil
}

func createAthlete(ctx *echo.Context, competitionID string, req *models.AthleteRequest) (*models.Athlete, error) {
	repository := repository.GetAthleteRepository()

	id, err := competitionIDFromPath(competitionID)
	if err != nil {
		return nil, err
	}

	athlete := models.Athlete{
		AthleteRequest: *req,
		CompetitionID:  id,
	}

	if err := repository.PutAthlete(ctx, &athlete); err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, exceptions.NewErrBadEntity(fmt.Sprintf("division with id %d does not exist", athlete.DivisionID))
		}
		return nil, err
	}

	logger.Action(ctx, "create_athlete", "competitionID", competitionID, "athleteID", athlete.ID)
	return &athlete, nil
}

func updateAthlete(ctx *echo.Context, competitionID string, athleteID string, req *models.AthleteRequest, athlete *models.Athlete) (*models.Athlete, error) {
	repository := repository.GetAthleteRepository()

	id, err := competitionIDFromPath(competitionID)
	if err != nil {
		return nil, err
	}

	if athlete == nil {
		a, err := repository.GetAthlete(ctx, athleteID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, exceptions.NewErrNotFound(fmt.Sprintf("athlete with id %s not found", athleteID))
			}
			return nil, err
		}
		athlete = a
	}
	if err := ensureAthleteInCompetition(athlete, competitionID); err != nil {
		return nil, err
	}
	athlete.AthleteRequest = *req
	athlete.CompetitionID = id

	if err := repository.PutAthlete(ctx, athlete); err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, exceptions.NewErrBadEntity(fmt.Sprintf("division with id %d does not exist", athlete.DivisionID))
		}
		return nil, err
	}

	logger.Action(ctx, "update_athlete", "competitionID", competitionID, "athleteID", athlete.ID)
	return athlete, nil
}

func deleteAthlete(ctx *echo.Context, competitionID string, id string) error {
	athleteRepository := repository.GetAthleteRepository()

	if _, err := getAthlete(ctx, competitionID, id); err != nil {
		return err
	}
	if err := athleteRepository.DeleteAthlete(ctx, id); err != nil {
		return err
	}

	logger.Action(ctx, "delete_athlete", "competitionID", competitionID, "athleteID", id)
	return nil
}

func ensureAthleteInCompetition(athlete *models.Athlete, competitionID string) error {
	if competitionID == "" {
		return nil
	}
	id, err := competitionIDFromPath(competitionID)
	if err != nil {
		return err
	}
	if athlete.CompetitionID != id {
		return exceptions.NewErrNotFound(fmt.Sprintf("athlete with id %d is not part of competition %s", athlete.ID, competitionID))
	}
	return nil
}

func competitionIDFromPath(competitionID string) (uint, error) {
	id, err := strconv.ParseUint(competitionID, 10, 64)
	if err != nil {
		return 0, exceptions.NewErrNotFound(fmt.Sprintf("competition with id %s not found", competitionID))
	}
	return uint(id), nil
}
