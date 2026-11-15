package models

import (
	"encoding/json"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type DivisionRequest struct {
	Name   string
	Gender string
}

type Division struct {
	gorm.Model
	DivisionRequest
}

func (d *Division) String() string {
	j, _ := json.Marshal(d)
	str, _ := strconv.Unquote(strings.ReplaceAll(strconv.Quote(string(j)), `\\u`, `\u`))

	return str
}
