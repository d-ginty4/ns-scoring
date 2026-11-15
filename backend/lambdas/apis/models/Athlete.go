package models

import (
	"encoding/json"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type AthleteRequest struct {
	Name    string
	Country string
}

type Athlete struct {
	gorm.Model
	AthleteRequest
}

func (a *Athlete) String() string {
	j, _ := json.Marshal(a)
	str, _ := strconv.Unquote(strings.ReplaceAll(strconv.Quote(string(j)), `\\u`, `\u`))

	return str
}
