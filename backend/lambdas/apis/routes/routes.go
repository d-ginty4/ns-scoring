package routes

import (
	"net/http"

	_ "apis/docs"
	"apis/middleware"

	"github.com/labstack/echo/v5"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

// @title APIs
// @version 1
// @description Documentation for all the apis. Every /api/v1 JSON response uses responses.Envelope with requestID, data, and error fields.
// @host localhost:3000
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT.
func SetupRoutes(e *echo.Echo) {
	regionRoutes := getRegionRoutes()
	divisionRoutes := getDivisionRoutes()
	competitionRoutes := getCompetitionRoutes()
	athleteRoutes := getAthleteRoutes()
	eventRoutes := getEventRoutes()
	resultRoutes := getResultRoutes()

	// Swagger docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Base path
	base := e.Group("/api/v1")

	// Region Routes (/api/v1/regions)
	regionGroup := base.Group("/regions")
	regionGroup.POST("", regionRoutes.PostRegion)
	regionGroup.GET("", regionRoutes.QueryRegions)
	regionGroupID := regionGroup.Group("/:regionID", middleware.CheckRegionExist)
	regionGroupID.GET("", regionRoutes.GetRegion)
	regionGroupID.PUT("", regionRoutes.PutRegion)
	regionGroupID.DELETE("", regionRoutes.DeleteRegion)

	// Division Routes (/api/v1/divisions)
	divisionGroup := base.Group("/divisions")
	divisionGroup.POST("", divisionRoutes.PostDivision)
	divisionGroup.GET("", divisionRoutes.QueryDivisions)
	divisionGroupID := divisionGroup.Group("/:divisionID", middleware.CheckDivisionExist)
	divisionGroupID.GET("", divisionRoutes.GetDivision)
	divisionGroupID.PUT("", divisionRoutes.PutDivision)
	divisionGroupID.DELETE("", divisionRoutes.DeleteDivision)

	// Competition Routes (/api/v1/competitions)
	competitionGroup := base.Group("/competitions")
	competitionGroup.POST("", competitionRoutes.PostCompetition)
	competitionGroup.GET("", competitionRoutes.QueryCompetitions)
	competitionGroupID := competitionGroup.Group("/:competitionID", middleware.CheckCompetitionExist)
	competitionGroupID.GET("", competitionRoutes.GetCompetition)
	competitionGroupID.PUT("", competitionRoutes.PutCompetition)
	competitionGroupID.DELETE("", competitionRoutes.DeleteCompetition)

	// Event Routes (/api/v1/events)
	eventGroup := base.Group("/events")
	eventGroup.POST("", eventRoutes.PostEvent)
	eventGroup.GET("", eventRoutes.QueryEvents)
	eventGroupID := eventGroup.Group("/:eventID", middleware.CheckEventExist)
	eventGroupID.GET("", eventRoutes.GetEvent)
	eventGroupID.PUT("", eventRoutes.PutEvent)
	eventGroupID.DELETE("", eventRoutes.DeleteEvent)

	// Athlete Routes (/api/v1/competitions/:competitionID/athletes)
	competitionGroupID.POST("/athletes", athleteRoutes.PostAthlete)
	competitionGroupID.POST("/athletes/bulk", athleteRoutes.PostAthletes)
	competitionGroupID.GET("/athletes", athleteRoutes.QueryAthletes)
	athleteGroupID := competitionGroupID.Group("/athletes/:athleteID", middleware.CheckAthleteExist)
	athleteGroupID.GET("", athleteRoutes.GetAthlete)
	athleteGroupID.PUT("", athleteRoutes.PutAthlete)
	athleteGroupID.DELETE("", athleteRoutes.DeleteAthlete)

	// Result Routes (/api/v1/results/athletes/:athleteID/events/:eventID)
	resultGroup := base.Group("/results")
	resultGroup.GET("/athletes/:athleteID", resultRoutes.QueryResults, middleware.CheckAthleteExist)
	resultGroup.POST("/athletes/:athleteID", resultRoutes.PostAthleteResults, middleware.CheckAthleteExist)
	resultEventGroup := resultGroup.Group("/athletes/:athleteID/events/:eventID", middleware.CheckAthleteExist, middleware.CheckEventExist)
	resultEventGroup.POST("", resultRoutes.PostResult)
	resultEventGroup.GET("", resultRoutes.GetResult, middleware.CheckResultExist)
	resultEventGroup.PUT("", resultRoutes.PutResult, middleware.CheckResultExist)
	resultEventGroup.DELETE("", resultRoutes.DeleteResult, middleware.CheckResultExist)

	// Special route to return all available routes
	e.GET("/routes", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, e.Router().Routes())
	})
}

// NotFound is the handler for requests that match no route.
func NotFound(c *echo.Context) error {
	return echo.NewHTTPError(http.StatusNotFound, "route not found")
}
