package service

import (
	"errors"
	"fmt"

	"apis/exceptions"
	"apis/logger"
	"apis/models"
	"apis/repository"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type EventService struct {
	QueryEvents func(ctx *echo.Context) (*[]models.Event, error)
	CreateEvent func(ctx *echo.Context, req *models.EventRequest) (*models.Event, error)
	GetEvent    func(ctx *echo.Context, id string) (*models.Event, error)
	UpdateEvent func(ctx *echo.Context, eventID string, req *models.EventRequest, existingEvent *models.Event) (*models.Event, error)
	DeleteEvent func(ctx *echo.Context, id string) error
}

var eventService *EventService

func GetEventService() *EventService {
	if eventService == nil {
		return newEventService()
	}
	return eventService
}

func newEventService() *EventService {
	service := &EventService{
		QueryEvents: queryEvents,
		CreateEvent: createEvent,
		GetEvent:    getEvent,
		UpdateEvent: updateEvent,
		DeleteEvent: deleteEvent,
	}

	eventService = service
	return service
}

func queryEvents(ctx *echo.Context) (*[]models.Event, error) {
	events, err := repository.GetEventRepository().QueryEvents(ctx)
	if err != nil {
		return nil, err
	}

	logger.Action(ctx, "query_events", "count", len(*events))
	return events, nil
}

func getEvent(ctx *echo.Context, id string) (*models.Event, error) {
	eventRepository := repository.GetEventRepository()

	event, err := eventRepository.GetEvent(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewErrNotFound(fmt.Sprintf("event with id %s not found", id))
		}
		return nil, err
	}

	logger.Action(ctx, "get_event", "eventID", event.ID)
	return event, nil
}

func createEvent(ctx *echo.Context, req *models.EventRequest) (*models.Event, error) {
	repository := repository.GetEventRepository()

	event := models.Event{
		EventRequest: *req,
	}

	if err := repository.PutEvent(ctx, &event); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, exceptions.NewErrBadEntity(fmt.Sprintf("event with name %s already exists", event.Name))
		}
		return nil, err
	}

	logger.Action(ctx, "create_event", "eventID", event.ID)
	return &event, nil
}

func updateEvent(ctx *echo.Context, eventID string, req *models.EventRequest, event *models.Event) (*models.Event, error) {
	repository := repository.GetEventRepository()

	if event == nil {
		e, err := repository.GetEvent(ctx, eventID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, exceptions.NewErrNotFound(fmt.Sprintf("event with id %s not found", eventID))
			}
			return nil, err
		}
		event = e
	}
	event.EventRequest = *req

	if err := repository.PutEvent(ctx, event); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, exceptions.NewErrBadEntity(fmt.Sprintf("event with name %s already exists", event.Name))
		}
		return nil, err
	}

	logger.Action(ctx, "update_event", "eventID", event.ID)
	return event, nil
}

func deleteEvent(ctx *echo.Context, id string) error {
	eventRepository := repository.GetEventRepository()

	if err := eventRepository.DeleteEvent(ctx, id); err != nil {
		return err
	}

	logger.Action(ctx, "delete_event", "eventID", id)
	return nil
}
