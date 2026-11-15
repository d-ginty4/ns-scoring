package models

import (
	"encoding/json"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type CompetitionRequest struct {
	Name     string
	Country  string
	IconName string
}

type Competition struct {
	gorm.Model
	CompetitionRequest
}

func (c *Competition) String() string {
	j, _ := json.Marshal(c)
	str, _ := strconv.Unquote(strings.ReplaceAll(strconv.Quote(string(j)), `\\u`, `\u`))

	return str
}
