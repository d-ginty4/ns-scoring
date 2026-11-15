package repository

import (
	"apis/models"

	"github.com/labstack/echo/v5"
)

type EventRepository struct {
	QueryEvents func(ctx *echo.Context) (*[]models.Event, error)
	PutEvent    func(ctx *echo.Context, event *models.Event) error
	GetEvent    func(ctx *echo.Context, eventID string) (*models.Event, error)
	DeleteEvent func(ctx *echo.Context, eventID string) error
}

var eventRepository *EventRepository

func GetEventRepository() *EventRepository {
	if eventRepository == nil {
		return newEventRepository()
	}

	return eventRepository
}

func newEventRepository() *EventRepository {
	repository := &EventRepository{
		QueryEvents: queryEvents,
		PutEvent:    putEvent,
		GetEvent:    getEvent,
		DeleteEvent: deleteEvent,
	}

	eventRepository = repository
	return repository
}

func queryEvents(ctx *echo.Context) (*[]models.Event, error) {
	return &[]models.Event{}, nil
}

func putEvent(ctx *echo.Context, event *models.Event) error {
	return nil
}

func getEvent(ctx *echo.Context, eventID string) (*models.Event, error) {
	return &models.Event{}, nil
}

func deleteEvent(ctx *echo.Context, eventID string) error {
	return nil
}
