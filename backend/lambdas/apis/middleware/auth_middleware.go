package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"apis/logger"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

var awsRegion = os.Getenv("AWS_REGION")
var cognitoPoolID = os.Getenv("COGNITO_POOL_ID")
var keysUrl = fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", awsRegion, cognitoPoolID)

func ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	if keysUrl == "" {
		panic("COGNITO_JWKS_URL is not set")
	}

	return func(ctx *echo.Context) error {
		log := logger.From(ctx)

		token := ctx.Request().Header.Get("Authorization")
		if token == "" {
			log.Warn("unauthorized", "reason", "missing_token")
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		jwks, err := keyfunc.NewDefault([]string{keysUrl})
		if err != nil {
			log.Error("failed to fetch jwks", logger.Err(err))
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		parts := strings.Split(token, " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			log.Warn("unauthorized", "reason", "malformed_authorization_header")
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		t, err := jwt.Parse(parts[1], jwks.Keyfunc)
		if err != nil {
			log.Warn("unauthorized", "reason", "invalid_token", logger.Err(err))
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		userID, err := t.Claims.GetSubject()
		if err != nil {
			log.Error("failed to get userID from token", logger.Err(err))
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		logger.SetActor(ctx, userID, "", false)
		return next(ctx)
	}
}
