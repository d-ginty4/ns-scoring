package models

import (
	"encoding/json"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type EventRequest struct {
	Name               string
	EventNumber        int
	WeightOptionMens   []int
	WeightOptionWomens []int
	ScoringMethod      string
}

type Event struct {
	gorm.Model
	EventRequest
}

func (e *Event) String() string {
	j, _ := json.Marshal(e)
	str, _ := strconv.Unquote(strings.ReplaceAll(strconv.Quote(string(j)), `\\u`, `\u`))

	return str
}
