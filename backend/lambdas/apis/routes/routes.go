package routes

import (
	"net/http"

	_ "apis/docs"
	"apis/middleware"

	"github.com/labstack/echo/v5"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

var middlewareGroups = middleware.GetMiddlewareGrous()

// @title APIs
// @version 0.1
// @description Documentation for all the apis. Every /api/v1 JSON response uses responses.Envelope with requestID, data, and error fields.
// @host localhost:3000
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT.
func SetupRoutes(e *echo.Echo) {
	regionRoutes := getRegionRoutes()

	// Swagger docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Base path
	base := e.Group("/api/v1")

	// Region Routes
	regionGroup := base.Group("/regions")
	regionGroup.POST("", regionRoutes.PostRegion)
	regionGroup.GET("", regionRoutes.QueryRegions)
	regionGroup.GET(":regionId", regionRoutes.GetRegion)
	regionGroup.PUT(":regionId", regionRoutes.PutRegion)
	regionGroup.DELETE(":regionId", regionRoutes.DeleteRegion)

	// Special route to return all available routes
	e.GET("/routes", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, e.Router().Routes())
	})
}
