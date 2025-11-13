package service

import (
	"context"

	"github.com/CogitoNTNU/cogi-go/internal/model"
	eventRepository "github.com/CogitoNTNU/cogi-go/internal/repository/event"
	"github.com/sirupsen/logrus"
)

type Event struct {
    logger *logrus.Entry
    repository *eventRepository.Repo
}

func NewEventService(eventRepository *eventRepository.Repo ) *Event {
    return &Event{repository: eventRepository}
}

func (e *Event) GetAllEvents(ctx *context.Context) ([]model.Event, *model.ErrorResponse) {
    allEvents, err := e.repository.GetAllEvents(ctx)
    
    if err != nil {
        e.logger.Errorf("An error has occurred when retrieving all events. Error code %s", err.Code)
        return nil, err
    }

    return allEvents, nil

}
