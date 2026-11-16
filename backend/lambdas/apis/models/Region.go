package models

import (
	"encoding/json"
	"strconv"
	"strings"
)

type RegionRequest struct {
	Name string `json:"name" gorm:"unique;not null" validate:"required"`
}

type Region struct {
	Common
	RegionRequest
}

func (r *Region) String() string {
	j, _ := json.Marshal(r)
	str, _ := strconv.Unquote(strings.ReplaceAll(strconv.Quote(string(j)), `\\u`, `\u`))

	return str
}
