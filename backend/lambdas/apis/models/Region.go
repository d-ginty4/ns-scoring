package models

import (
	"encoding/json"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type RegionRequest struct {
	Name string
}

type Region struct {
	gorm.Model
	RegionRequest
}

func (r *Region) String() string {
	j, _ := json.Marshal(r)
	str, _ := strconv.Unquote(strings.ReplaceAll(strconv.Quote(string(j)), `\\u`, `\u`))

	return str
}
