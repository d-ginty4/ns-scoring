package repository

import (
	"apis/models"

	"github.com/labstack/echo/v5"
)

type RegionRepository struct {
	QueryRegions func(ctx *echo.Context) (*[]models.Region, error)
	PutRegion    func(ctx *echo.Context, region *models.Region) error
	GetRegion    func(ctx *echo.Context, regionID string) (*models.Region, error)
	DeleteRegion func(ctx *echo.Context, regionID string) error
}

var regionRepository *RegionRepository

func GetRegionRepository() *RegionRepository {
	if regionRepository == nil {
		return newRegionRepository()
	}

	return regionRepository
}

func newRegionRepository() *RegionRepository {
	repository := &RegionRepository{
		QueryRegions: queryRegions,
		PutRegion:    putRegion,
		GetRegion:    getRegion,
		DeleteRegion: deleteRegion,
	}

	regionRepository = repository
	return repository
}

func queryRegions(ctx *echo.Context) (*[]models.Region, error) {
	if db == nil {
		ConnectDB()
	}
	var regions []models.Region
	return &regions, (db.Find(&regions)).Error
}

func getRegion(ctx *echo.Context, regionID string) (*models.Region, error) {
	if db == nil {
		ConnectDB()
	}

	var region models.Region
	result := db.First(&region, regionID)
	if result.Error != nil {
		return nil, result.Error
	}

	return &region, nil
}

func putRegion(ctx *echo.Context, region *models.Region) error {
	if db == nil {
		ConnectDB()
	}
	return db.Save(region).Error
}

func deleteRegion(ctx *echo.Context, regionID string) error {
	if db == nil {
		ConnectDB()
	}
	return (db.Unscoped().Delete(&models.Region{}, regionID)).Error
}
