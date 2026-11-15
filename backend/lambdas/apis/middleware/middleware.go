package middleware

import "github.com/labstack/echo/v5"

type MiddlewareGroups struct {
	// UserAuthGroup will call the middleware ValidateToken and ExtractUser
	UserAuthGroup func(next echo.HandlerFunc) echo.HandlerFunc
}

func GetMiddlewareGrous() *MiddlewareGroups {
	return &MiddlewareGroups{}
}

// func userAuthGroup(next echo.HandlerFunc) echo.HandlerFunc {
// 	return ValidateToken(ExtractUser(next))
// }
