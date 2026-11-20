package utils

import (
	"apis/exceptions"

	"github.com/labstack/echo/v5"
)

func BindAndValidate(ctx *echo.Context, target any) error {
	if err := ctx.Bind(target); err != nil {
		return &exceptions.ErrValidation{Errors: []string{err.Error()}}
	}
	return ctx.Validate(target)
}
