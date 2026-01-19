package tempApplicationRepository

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/CogitoNTNU/cogi-go/internal/api/dto"
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

func (r *Repo) GetAllTempApplications(ctx *context.Context) ([]model.TempApplication, *model.ErrorResponse) {
	cCtx, cancel := context.WithTimeout(*ctx, r.queryTimeoutLimit)
	defer cancel()

	var tempApplicationResults []db.TempApplication
	args := map[string]any{}

	err := r.queries.Read.getAllTempApplications.SelectContext(cCtx, &tempApplicationResults, args)
	if err != nil {
		return nil, &model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	tempApplications := make([]model.TempApplication, 0, len(tempApplicationResults))
	for _, tempApplicationResult := range tempApplicationResults {
		tempApplication := tempApplicationResult.ToModel()
		tempApplications = append(tempApplications, *tempApplication)
	}

	return tempApplications, nil
}

func (r *Repo) InsertTempApplication(ctx *context.Context, tempApplication *dto.CreateTempApplicationRequest) *model.ErrorResponse {
	cCtx, cancel := context.WithTimeout(*ctx, r.queryTimeoutLimit)
	defer cancel()

	createdAt := time.Now()

	dbTempApplication := db.TempApplication{
		Id:              uuid.New(),
		FirstName:       tempApplication.FirstName,
		LastName:        tempApplication.LastName,
		Email:           tempApplication.Email,
		PhoneNumber:     tempApplication.PhoneNumber,
		Projects:        tempApplication.Projects,
		ApplicationText: tempApplication.ApplicationText,
		CreatedAt:       &createdAt,
	}

	_, err := r.queries.Write.insertTempApplication.ExecContext(cCtx, dbTempApplication)
	if err != nil {
		return &model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

func (r *Repo) Close() (err error) {
	err = r.queries.Close()
	if err != nil {
		return fmt.Errorf("error closing queries %w", err)
	}
	return nil
}
