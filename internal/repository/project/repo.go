package projectRepository

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

func (r *Repo) GetAllProjects(ctx *context.Context) ([]model.Project, *model.ErrorResponse) {
	cCtx, cancel := context.WithTimeout(*ctx, r.queryTimeoutLimit)
	defer cancel()

	var projectResults []db.Project
	args := map[string]any{}

	err := r.queries.Read.getProjects.SelectContext(cCtx, &projectResults, args)
	if err != nil {
		return nil, &model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	projects := make([]model.Project, 0, len(projectResults))
	for _, projectResult := range projectResults {
		project := projectResult.ToModel()
		projects = append(projects, *project)
	}

	return projects, nil
}

func (r *Repo) GetProjectByID(ctx *context.Context, projectId uuid.UUID) (*model.Project, *model.ErrorResponse) {
	cCtx, cancel := context.WithTimeout(*ctx, r.queryTimeoutLimit)
	defer cancel()

	var projectResult db.Project
	args := map[string]any{
		"projectId": projectId,
	}

	err := r.queries.Read.getProject.GetContext(cCtx, &projectResult, args)
	if err != nil {
		return nil, &model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	project := projectResult.ToModel()

	return project, nil
}

func (r *Repo) Close() (err error) {
	err = r.queries.Close()
	if err != nil {
		return fmt.Errorf("error closing queries %w", err)
	}
	return
}
