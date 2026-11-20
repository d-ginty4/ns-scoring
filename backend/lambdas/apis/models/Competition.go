package models

import (
	"encoding/json"
	"strconv"
	"strings"
)

type CompetitionRequest struct {
	Name     string `json:"name" gorm:"not null" validate:"required"`
	Country  string `json:"country" gorm:"not null" validate:"required"`
	IconName string `json:"iconName"`
	RegionID uint   `json:"regionId" gorm:"not null;index" validate:"required,regionExists"`
}

type Competition struct {
	Common
	CompetitionRequest
	Region Region `json:"region" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (c *Competition) String() string {
	j, _ := json.Marshal(c)
	str, _ := strconv.Unquote(strings.ReplaceAll(strconv.Quote(string(j)), `\\u`, `\u`))

	return str
}
