package models

import (
	"encoding/json"
	"strconv"
	"strings"
)

type EventRequest struct {
	Name               string `json:"name" gorm:"unique;not null" validate:"required"`
	EventNumber        int    `json:"eventNumber" gorm:"not null" validate:"required"`
	WeightOptionsMen   []int  `json:"weightOptionsMen" gorm:"type:jsonb;serializer:json;not null" validate:"required"`
	WeightOptionsWomen []int  `json:"weightOptionsWomen" gorm:"type:jsonb;serializer:json;not null" validate:"required"`
	ScoringMethod      string `json:"scoringMethod" gorm:"not null" validate:"required"`
}

type Event struct {
	Common
	EventRequest
}

func (e *Event) String() string {
	j, _ := json.Marshal(e)
	str, _ := strconv.Unquote(strings.ReplaceAll(strconv.Quote(string(j)), `\\u`, `\u`))

	return str
}
