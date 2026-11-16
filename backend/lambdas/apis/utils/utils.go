package utils

import (
	"strings"
	"unicode"

	"apis/exceptions"

	"github.com/labstack/echo/v5"
)

func BindAndValidate(ctx *echo.Context, target any) error {
	if err := ctx.Bind(target); err != nil {
		return &exceptions.ErrValidation{Errors: []string{err.Error()}}
	}
	return ctx.Validate(target)
}

func ToCamelCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})

	if len(words) == 0 {
		return ""
	}

	var result strings.Builder

	for i, word := range words {
		word = strings.ToLower(word)

		if i > 0 {
			word = strings.ToUpper(word[:1]) + word[1:]
		}

		result.WriteString(word)
	}

	return result.String()
}
