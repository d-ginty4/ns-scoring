package routes

import (
	"errors"
	"net/http"

	"apis/exceptions"
	"apis/models"
	"apis/responses"
	"apis/service"
	"apis/utils"

	"github.com/labstack/echo/v5"
)

type EventRoutes struct {
	QueryEvents func(c *echo.Context) error
	PostEvent   func(c *echo.Context) error
	GetEvent    func(c *echo.Context) error
	PutEvent    func(c *echo.Context) error
	DeleteEvent func(c *echo.Context) error
}

var eventRoutes *EventRoutes

func getEventRoutes() *EventRoutes {
	if eventRoutes == nil {
		return newEventRoutes()
	}
	return eventRoutes
}

func newEventRoutes() *EventRoutes {
	routes := &EventRoutes{
		QueryEvents: queryEvents,
		PostEvent:   postEvent,
		GetEvent:    getEvent,
		PutEvent:    putEvent,
		DeleteEvent: deleteEvent,
	}

	eventRoutes = routes
	return routes
}

func queryEvents(ctx *echo.Context) error {
	service := service.GetEventService()

	page, err := service.QueryEvents(ctx)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, page)
}

func postEvent(ctx *echo.Context) error {
	service := service.GetEventService()

	requestUser, ok := ctx.Get("requestUser").(*models.User)
	if !ok || requestUser == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get request user")
	}

	var p models.EventRequest
	if err := utils.BindAndValidate(ctx, &p); err != nil {
		return err
	}

	event, err := service.CreateEvent(ctx, &p, requestUser)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, event)
}

func getEvent(ctx *echo.Context) error {
	eventService := service.GetEventService()

	id := ctx.Param("eventID")

	event, ok := ctx.Get("event").(*models.Event)
	if ok && event != nil {
		return responses.Success(ctx, http.StatusOK, event)
	}

	event, err := eventService.GetEvent(ctx, id)
	if err != nil {
		if errors.Is(err, &exceptions.ErrNotFound{}) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return responses.Success(ctx, http.StatusOK, event)
}

func putEvent(ctx *echo.Context) error {
	service := service.GetEventService()

	existing, ok := ctx.Get("event").(*models.Event)
	if !ok || existing == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get existing event")
	}

	var p models.EventRequest
	if err := utils.BindAndValidate(ctx, &p); err != nil {
		return err
	}

	id := ctx.Param("eventID")
	event, err := service.UpdateEvent(ctx, id, &p, existing)
	if err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, event)
}

func deleteEvent(ctx *echo.Context) error {
	service := service.GetEventService()

	event, ok := ctx.Get("event").(*models.Event)
	if !ok || event == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get event")
	}

	id := ctx.Param("eventID")
	if err := service.DeleteEvent(ctx, id, event); err != nil {
		return err
	}

	return responses.Success(ctx, http.StatusOK, nil)
}
