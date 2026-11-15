package repository

import (
	"apis/models"

	"github.com/labstack/echo/v5"
)

type RegionRepository struct {
	QueryRegions func(ctx *echo.Context) (*[]models.Region, error)
	CreateRegion func(ctx *echo.Context, region *models.Region) error
	UpdateRegion func(ctx *echo.Context, region *models.Region) error
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
		CreateRegion: createRegion,
		UpdateRegion: updateRegion,
		GetRegion:    getRegion,
		DeleteRegion: deleteRegion,
	}

	regionRepository = repository
	return repository
}

func queryRegions(ctx *echo.Context) (*[]models.Region, error) {
	if db == nil {
		connectDB()
	}
	var regions []models.Region
	return &regions, (db.Find(&regions)).Error
}

func createRegion(ctx *echo.Context, region *models.Region) error {
	if db == nil {
		connectDB()
	}
	return (db.Create(region)).Error
}

func updateRegion(ctx *echo.Context, region *models.Region) error {
	if db == nil {
		connectDB()
	}
	return nil
}

func getRegion(ctx *echo.Context, regionID string) (*models.Region, error) {
	if db == nil {
		connectDB()
	}
	var region models.Region
	return &region, (db.First(&region, regionID)).Error
}

func deleteRegion(ctx *echo.Context, regionID string) error {
	if db == nil {
		connectDB()
	}
	return (db.Delete(&models.Region{}, regionID)).Error
}
