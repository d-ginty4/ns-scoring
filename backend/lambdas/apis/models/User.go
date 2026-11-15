package models

import (
	"encoding/json"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type UserRequest struct {
	Name    string
	IsAdmin bool
	email   string
}

type User struct {
	gorm.Model
	UserRequest
}

func (u *User) String() string {
	j, _ := json.Marshal(u)
	str, _ := strconv.Unquote(strings.ReplaceAll(strconv.Quote(string(j)), `\\u`, `\u`))

	return str
}
