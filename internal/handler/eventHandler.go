package handler

import (
	"context"
	"net/http"

	"github.com/CogitoNTNU/cogi-go/internal/service"
	"github.com/gin-gonic/gin"
)

type Event struct {
    service *service.Event
    ctx *context.Context
}

func NewEventHandler(eventService *service.Event, ctx *context.Context) *Event {
    return &Event{service: eventService, ctx: ctx}
}

func (e *Event) GetAllEvents(gCtx *gin.Context) {
    allEvents, err := e.service.GetAllEvents(e.ctx)
    if err != nil {
        gCtx.AbortWithStatusJSON(err.Code, err.Message)
        return
    }

    gCtx.JSON(http.StatusOK, allEvents)
}

func (e *Event) GetEventByID(gCtx *gin.Context) {
    eventId := gCtx.Param("eventId")
    event, err := e.service.GetEventByID(e.ctx, eventId)
    
    if err != nil {
        gCtx.AbortWithStatusJSON(err.Code, err.Message)
        return
    }

    gCtx.JSON(http.StatusOK, event)
}