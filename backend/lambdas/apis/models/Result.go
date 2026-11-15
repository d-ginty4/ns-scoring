package models

import (
	"encoding/json"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type ResultRequest struct {
	Weight            int
	Value             int
	DidNotParticipate bool
}

type Result struct {
	gorm.Model
	ResultRequest
}

func (r *Result) String() string {
	j, _ := json.Marshal(r)
	str, _ := strconv.Unquote(strings.ReplaceAll(strconv.Quote(string(j)), `\\u`, `\u`))

	return str
}
