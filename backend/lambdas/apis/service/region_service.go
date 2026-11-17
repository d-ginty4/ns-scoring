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

type RegionService struct {
	QueryRegions func(ctx *echo.Context) (*[]models.Region, error)
	CreateRegion func(ctx *echo.Context, req *models.RegionRequest) (*models.Region, error)
	GetRegion    func(ctx *echo.Context, id string) (*models.Region, error)
	UpdateRegion func(ctx *echo.Context, regionID string, req *models.RegionRequest, existingRegion *models.Region) (*models.Region, error)
	DeleteRegion func(ctx *echo.Context, id string) error
}

var regionService *RegionService

func GetRegionService() *RegionService {
	if regionService == nil {
		return newRegionService()
	}
	return regionService
}

func newRegionService() *RegionService {
	service := &RegionService{
		QueryRegions: queryRegions,
		CreateRegion: createRegion,
		GetRegion:    getRegion,
		UpdateRegion: updateRegion,
		DeleteRegion: deleteRegion,
	}

	regionService = service
	return service
}

func queryRegions(ctx *echo.Context) (*[]models.Region, error) {
	regions, err := repository.GetRegionRepository().QueryRegions(ctx)
	if err != nil {
		return nil, err
	}

	logger.Action(ctx, "query_regions", "count", len(*regions))
	return regions, nil
}

func getRegion(ctx *echo.Context, id string) (*models.Region, error) {
	regionRepository := repository.GetRegionRepository()

	region, err := regionRepository.GetRegion(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewErrNotFound(fmt.Sprintf("region with id %s not found", id))
		}
		return nil, err
	}

	logger.Action(ctx, "get_region", "regionID", region.ID)
	return region, nil
}

func createRegion(ctx *echo.Context, req *models.RegionRequest) (*models.Region, error) {
	repository := repository.GetRegionRepository()

	region := models.Region{
		RegionRequest: *req,
	}

	if err := repository.PutRegion(ctx, &region); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, exceptions.NewErrBadEntity(fmt.Sprintf("region with name %s already exists", region.Name))
		}
		return nil, err
	}

	logger.Action(ctx, "create_region", "regionID", region.ID)
	return &region, nil
}

func updateRegion(ctx *echo.Context, regionID string, req *models.RegionRequest, region *models.Region) (*models.Region, error) {
	repository := repository.GetRegionRepository()

	if region == nil {
		r, err := repository.GetRegion(ctx, regionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, exceptions.NewErrNotFound(fmt.Sprintf("region with id %s not found", regionID))
			}
			return nil, err
		}
		region = r
	}
	region.RegionRequest = *req

	if err := repository.PutRegion(ctx, region); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, exceptions.NewErrBadEntity(fmt.Sprintf("region with name %s already exists", region.Name))
		}
		return nil, err
	}

	logger.Action(ctx, "update_region", "regionID", region.ID)
	return region, nil
}

func deleteRegion(ctx *echo.Context, id string) error {
	regionRepository := repository.GetRegionRepository()

	if err := regionRepository.DeleteRegion(ctx, id); err != nil {
		return err
	}

	logger.Action(ctx, "delete_region", "regionID", id)
	return nil
}
