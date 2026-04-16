package handler

import (
	"context"
	"net/http"

	"github.com/CogitoNTNU/cogi-go/internal/api/dto"
	"github.com/CogitoNTNU/cogi-go/internal/service"
	"github.com/gin-gonic/gin"
)

type Event struct {
	service *service.Event
	ctx     *context.Context
}

func NewEventHandler(eventService *service.Event, ctx *context.Context) *Event {
	return &Event{service: eventService, ctx: ctx}
}

func (e *Event) GetAllEvents(gCtx *gin.Context) {
	events, err := e.service.GetAllEvents(e.ctx)
	if err != nil {
		gCtx.AbortWithStatusJSON(err.Code, err.Message)
		return
	}

	gCtx.JSON(http.StatusOK, events)
}

func (e *Event) GetEventById(gCtx *gin.Context) {
	id := gCtx.Param("eventId")

	event, err := e.service.GetEventByID(e.ctx, id)
	if err != nil {
		gCtx.AbortWithStatusJSON(err.Code, err.Message)
		return
	}

	gCtx.JSON(http.StatusOK, event)
}

func (e *Event) CreateEvent(gCtx *gin.Context) {
	var req dto.CreateEventRequest

	err := gCtx.ShouldBindJSON(&req)
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err2 := e.service.CreateEvent(e.ctx, &req)
	if err2 != nil {
		gCtx.AbortWithStatusJSON(err2.Code, err2.Message)
		return
	}

	gCtx.Status(http.StatusCreated)
}
