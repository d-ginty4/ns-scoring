package models

import (
	"encoding/json"
	"strconv"
	"strings"
)

type ResultRequest struct {
	Weight            int  `json:"weight" gorm:"not null" validate:"required"`
	Value             int  `json:"value" gorm:"not null" validate:"required"`
	DidNotParticipate bool `json:"didNotParticipate"`
}

type AthleteResultItem struct {
	EventID uint `json:"eventId" validate:"required,eventExists"`
	ResultRequest
}

type AthleteResultsRequest struct {
	Results []AthleteResultItem `json:"results" validate:"required,min=1,dive"`
}

type Result struct {
	Common
	ResultRequest
	EventId   uint    `json:"eventId" gorm:"not null;index;uniqueIndex:idx_results_athlete_event,priority:2"`
	Event     Event   `json:"event" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AthleteId uint    `json:"athleteId" gorm:"not null;index;uniqueIndex:idx_results_athlete_event,priority:1"`
	Athlete   Athlete `json:"athlete" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (r *Result) String() string {
	j, _ := json.Marshal(r)
	str, _ := strconv.Unquote(strings.ReplaceAll(strconv.Quote(string(j)), `\\u`, `\u`))

	return str
}
