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

type Result struct {
	Common
	ResultRequest
	EventId   uint    `json:"eventId" gorm:"not null;index" validate:"required,eventExists"`
	Event     Event   `json:"event" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AthleteId uint    `json:"athleteId" gorm:"not null;index" validate:"required,athleteExists"`
	Athlete   Athlete `json:"athlete" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (r *Result) String() string {
	j, _ := json.Marshal(r)
	str, _ := strconv.Unquote(strings.ReplaceAll(strconv.Quote(string(j)), `\\u`, `\u`))

	return str
}
