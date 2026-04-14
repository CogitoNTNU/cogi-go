package eventRepository

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/CogitoNTNU/cogi-go/internal/model"
	"github.com/CogitoNTNU/cogi-go/internal/repository/db"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Repo struct {
	db                *sqlx.DB
	queries           queries
	queryTimeOutLimit time.Duration
	log               logrus.FieldLogger
}

func NewRepo(db *sqlx.DB, queryTimeOutLimit time.Duration, logger *logrus.Entry) *Repo {
	queries, err := PreparedQueries(db)

	if err != nil {
		logger.Fatalf("failed to prepare queries: %s", err)
	}

	return &Repo{
		db:                db,
		queries:           queries,
		queryTimeOutLimit: queryTimeOutLimit,
		log:               logger,
	}
}

func (r *Repo) GetAllEvents(ctx *context.Context) ([]model.Event, *model.ErrorResponse) {
	cCtx, cancel := context.WithTimeout(*ctx, r.queryTimeOutLimit)
	defer cancel()

	var eventResults []db.Event

	args := map[string]any{}
	err := r.queries.Read.getEvents.SelectContext(cCtx, &eventResults, args)

	if err != nil {
		return nil, &model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	events := make([]model.Event, 0, len(eventResults))
	for _, eventResult := range eventResults {
		event := eventResult.ToModel()
		events = append(events, *event)
	}

	return events, nil
}

func (r *Repo) GetEventByID(ctx *context.Context, eventId uuid.UUID) (*model.Event, *model.ErrorResponse) {
	cCtx, cancel := context.WithTimeout(*ctx, r.queryTimeOutLimit)
	defer cancel()

	var eventResult db.Event
	args := map[string]any{"eventId": eventId}
	err := r.queries.Read.getEvent.SelectContext(cCtx, &eventResult, args)

	if err != nil {
		return nil, &model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	event := eventResult.ToModel()

	return event, nil
}

func (r *Repo) Close() (err error) {
	err = r.queries.Close()
	if err != nil {
		return fmt.Errorf("error closing queries: %w", err)
	}
	return
}
