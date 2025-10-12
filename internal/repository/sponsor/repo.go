package sponsorRepository

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
	queryTimeoutLimit time.Duration
	log               logrus.FieldLogger
}

func NewRepo(db *sqlx.DB, queryTimeoutLimit time.Duration, logger *logrus.Entry) *Repo {
	queries, err := PrepareQueries(db)

	if err != nil {
		logger.Fatalf("Failed to prepare queries: %s", err)
	}

	return &Repo{
		db:                db,
		queries:           queries,
		queryTimeoutLimit: queryTimeoutLimit,
		log:               logger,
	}
}

func (r *Repo) GetAllSponsors(ctx *context.Context) ([]model.Sponsor, *model.ErrorResponse) {
	cCtx, cancel := context.WithTimeout(*ctx, r.queryTimeoutLimit)
	defer cancel()

	var sponsorResults []db.Sponsor
	args := map[string]any{}

	err := r.queries.Read.getSponsors.SelectContext(cCtx, &sponsorResults, args)
	if err != nil {
		return nil, &model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	sponsors := make([]model.Sponsor, 0, len(sponsorResults))
	for _, sponsorResult := range sponsorResults {
		sponsor := sponsorResult.ToModel()
		sponsors = append(sponsors, *sponsor)
	}

	return sponsors, nil
}

func (r *Repo) GetSponsorByID(ctx *context.Context, sponsorID uuid.UUID) (*model.Sponsor, *model.ErrorResponse) {
	cCtx, cancel := context.WithTimeout(*ctx, r.queryTimeoutLimit)
	defer cancel()

	var sponsorResult db.Sponsor
	args := map[string]any{
		"sponsor_id": sponsorID,
	}

	err := r.queries.Read.getSponsor.GetContext(cCtx, &sponsorResult, args)
	if err != nil {
		return nil, &model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	sponsor := sponsorResult.ToModel()
	return sponsor, nil
}

func (r *Repo) Close() (err error) {
	err = r.queries.Close()
	if err != nil {
		return fmt.Errorf("error closing queries %w", err)
	}
	
	return
}