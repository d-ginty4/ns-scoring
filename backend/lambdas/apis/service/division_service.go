package service

import (
	"errors"
	"fmt"

	"apis/exceptions"
	"apis/logger"
	"apis/models"
	"apis/repository"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type DivisionService struct {
	QueryDivisions func(ctx *echo.Context) (*[]models.Division, error)
	CreateDivision func(ctx *echo.Context, req *models.DivisionRequest) (*models.Division, error)
	GetDivision    func(ctx *echo.Context, id string) (*models.Division, error)
	UpdateDivision func(ctx *echo.Context, divisionID string, req *models.DivisionRequest, existingDivision *models.Division) (*models.Division, error)
	DeleteDivision func(ctx *echo.Context, id string) error
}

var divisionService *DivisionService

func GetDivisionService() *DivisionService {
	if divisionService == nil {
		return newDivisionService()
	}
	return divisionService
}

func newDivisionService() *DivisionService {
	service := &DivisionService{
		QueryDivisions: queryDivisions,
		CreateDivision: createDivision,
		GetDivision:    getDivision,
		UpdateDivision: updateDivision,
		DeleteDivision: deleteDivision,
	}

	divisionService = service
	return service
}

func queryDivisions(ctx *echo.Context) (*[]models.Division, error) {
	divisions, err := repository.GetDivisionRepository().QueryDivisions(ctx)
	if err != nil {
		return nil, err
	}

	logger.Action(ctx, "query_divisions", "count", len(*divisions))
	return divisions, nil
}

func getDivision(ctx *echo.Context, id string) (*models.Division, error) {
	divisionRepository := repository.GetDivisionRepository()

	division, err := divisionRepository.GetDivision(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewErrNotFound(fmt.Sprintf("division with id %s not found", id))
		}
		return nil, err
	}

	logger.Action(ctx, "get_division", "divisionID", division.ID)
	return division, nil
}

func createDivision(ctx *echo.Context, req *models.DivisionRequest) (*models.Division, error) {
	repository := repository.GetDivisionRepository()

	division := models.Division{
		DivisionRequest: *req,
	}

	if err := repository.PutDivision(ctx, &division); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, exceptions.NewErrBadEntity(fmt.Sprintf("%s division with name %s already exists", division.Gender, division.Name))
		}
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			return nil, exceptions.NewErrBadEntity("gender must be male or female")
		}
		return nil, err
	}

	logger.Action(ctx, "create_division", "divisionID", division.ID)
	return &division, nil
}

func updateDivision(ctx *echo.Context, divisionID string, req *models.DivisionRequest, division *models.Division) (*models.Division, error) {
	repository := repository.GetDivisionRepository()

	if division == nil {
		d, err := repository.GetDivision(ctx, divisionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, exceptions.NewErrNotFound(fmt.Sprintf("division with id %s not found", divisionID))
			}
			return nil, err
		}
		division = d
	}
	division.DivisionRequest = *req

	if err := repository.PutDivision(ctx, division); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, exceptions.NewErrBadEntity(fmt.Sprintf("%s division with name %s already exists", division.Gender, division.Name))
		}
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			return nil, exceptions.NewErrBadEntity("gender must be male or female")
		}
		return nil, err
	}

	logger.Action(ctx, "update_division", "divisionID", division.ID)
	return division, nil
}

func deleteDivision(ctx *echo.Context, id string) error {
	divisionRepository := repository.GetDivisionRepository()

	if err := divisionRepository.DeleteDivision(ctx, id); err != nil {
		return err
	}

	logger.Action(ctx, "delete_division", "divisionID", id)
	return nil
}
