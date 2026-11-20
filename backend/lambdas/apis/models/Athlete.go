package models

import (
	"encoding/json"
	"strconv"
	"strings"
)

type AthleteRequest struct {
	Name       string `json:"name" gorm:"not null" validate:"required"`
	Country    string `json:"country" gorm:"not null" validate:"required"`
	DivisionID uint   `json:"divisionId" gorm:"not null;index" validate:"required,divisionExists"`
}

type Athlete struct {
	Common
	AthleteRequest
	CompetitionID uint        `json:"competitionId" gorm:"not null;index" validate:"required,competitionExists"`
	Competition   Competition `json:"competition" gorm:"foreignKey:CompetitionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Division      Division    `json:"division" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (a *Athlete) String() string {
	j, _ := json.Marshal(a)
	str, _ := strconv.Unquote(strings.ReplaceAll(strconv.Quote(string(j)), `\\u`, `\u`))

	return str
}
