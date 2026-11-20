package validators

import (
	"errors"
	"strconv"

	"apis/repository"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func RegionExists(fl validator.FieldLevel) bool {
	return recordExists(fl, func(id string) error {
		_, err := repository.GetRegionRepository().GetRegion(nil, id)
		return err
	})
}

func DivisionExists(fl validator.FieldLevel) bool {
	return recordExists(fl, func(id string) error {
		_, err := repository.GetDivisionRepository().GetDivision(nil, id)
		return err
	})
}

func CompetitionExists(fl validator.FieldLevel) bool {
	return recordExists(fl, func(id string) error {
		_, err := repository.GetCompetitionRepository().GetCompetition(nil, id)
		return err
	})
}

func EventExists(fl validator.FieldLevel) bool {
	return recordExists(fl, func(id string) error {
		_, err := repository.GetEventRepository().GetEvent(nil, id)
		return err
	})
}

func AthleteExists(fl validator.FieldLevel) bool {
	return recordExists(fl, func(id string) error {
		_, err := repository.GetAthleteRepository().GetAthlete(nil, id)
		return err
	})
}

func recordExists(fl validator.FieldLevel, lookup func(id string) error) bool {
	id := fl.Field().Uint()
	if id == 0 {
		return true
	}

	err := lookup(strconv.FormatUint(id, 10))
	if err == nil {
		return true
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return false
}
