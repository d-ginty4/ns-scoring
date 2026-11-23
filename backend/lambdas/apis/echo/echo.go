package echo

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"apis/exceptions"
	"apis/logger"
	"apis/middleware"
	"apis/responses"
	"apis/routes"
	"apis/validators"

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
		errors = append(errors, fmt.Sprintf("'%s' %s", fieldLabel(fe), describe(fe)))
	}

	return &exceptions.ErrValidation{Errors: errors}
}

// fieldLabel keeps a top-level failure as the field name and includes the
// index for nested items, for example athletes[1].divisionId.
func fieldLabel(fe validator.FieldError) string {
	ns := fe.Namespace()
	if i := strings.Index(ns, "."); i >= 0 && i < len(ns)-1 {
		return ns[i+1:]
	}
	return fe.Field()
}

func describe(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "oneof":
		return fmt.Sprintf("must be one of: [%s]", fe.Param())
	case "min":
		if fe.Kind() == reflect.Slice || fe.Kind() == reflect.Array || fe.Kind() == reflect.Map {
			return "must contain at least " + fe.Param() + " items"
		}
		return "must be at least " + fe.Param() + " characters"
	case "max":
		return "must be at most " + fe.Param() + " characters"
	case "regionExists", "divisionExists", "competitionExists", "eventExists", "athleteExists":
		return fmt.Sprintf("%v does not exist", fe.Value())
	default:
		return "is invalid"
	}
}

func NewEchoInstance() *echo.Echo {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	v.RegisterValidation("regionExists", validators.RegionExists)
	v.RegisterValidation("divisionExists", validators.DivisionExists)
	v.RegisterValidation("competitionExists", validators.CompetitionExists)
	v.RegisterValidation("eventExists", validators.EventExists)
	v.RegisterValidation("athleteExists", validators.AthleteExists)

	e := echo.NewWithConfig(echo.Config{
		Router: echo.NewRouter(echo.RouterConfig{
			AllowOverwritingRoute: true,
			// Group middleware registers its own not-found routes with this handler,
			// so unknown paths under those groups use the same response.
			NotFoundHandler: routes.NotFound,
		}),
	})
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
