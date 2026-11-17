package models

import (
	"encoding/json"
	"strconv"
	"strings"
)

type DivisionRequest struct {
	Name   string `json:"name" gorm:"not null;uniqueIndex:idx_divisions_name_gender" validate:"required"`
	Gender string `json:"gender" gorm:"not null;uniqueIndex:idx_divisions_name_gender;check:chk_divisions_gender,gender IN ('male', 'female')" validate:"required,oneof=male female"`
}

type Division struct {
	Common
	DivisionRequest
}

func (d *Division) String() string {
	j, _ := json.Marshal(d)
	str, _ := strconv.Unquote(strings.ReplaceAll(strconv.Quote(string(j)), `\\u`, `\u`))

	return str
}
