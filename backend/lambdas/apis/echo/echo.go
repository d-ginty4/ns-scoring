package echo

import (
	"errors"
	"fmt"

	"apis/exceptions"
	"apis/logger"
	"apis/middleware"
	"apis/responses"
	"apis/routes"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
)

type structValidator struct {
	validator *validator.Validate
}

func (sv *structValidator) Validate(i any) error {
	err := sv.validator.Struct(i)
	if err == nil {
		return nil
	}

	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return err
	}

	errors := make([]string, 0)
	for _, fe := range verrs {
		errors = append(errors, fmt.Sprintf("'%s' %s", fe.Field(), describe(fe)))
	}

	return &exceptions.ErrValidation{Errors: errors}
}

func describe(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "oneof":
		return fmt.Sprintf("must be one of: [%s]", fe.Param())
	case "min":
		return "must be at least " + fe.Param() + " characters"
	case "max":
		return "must be at most " + fe.Param() + " characters"
	case "isCursor":
		return "invalid cursor"
	default:
		return "is invalid"
	}
}

func NewEchoInstance() *echo.Echo {
	v := validator.New(validator.WithRequiredStructEnabled())

	e := echo.New()
	e.Logger = logger.Logger()
	e.Validator = &structValidator{validator: v}
	e.HTTPErrorHandler = responses.HTTPErrorHandler

	// AccessLog is outermost so it still records panics recovered below.
	e.Use(middleware.AccessLog())
	e.Use(echoMiddleware.Recover())
	e.Use(middleware.RequestID())

	routes.SetupRoutes(e)
	return e
}
