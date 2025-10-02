package userRepository

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

func (r *Repo) GetAllUsers(ctx *context.Context) ([]model.User, *model.ErrorResponse) {
	cCtx, cancel := context.WithTimeout(*ctx, r.queryTimeoutLimit)
	defer cancel()

	var userResults []db.User
	args := map[string]any{}

	err := r.queries.Read.getUsers.SelectContext(cCtx, &userResults, args)
	if err != nil {
		return nil, &model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	users := make([]model.User, 0, len(userResults))
	for _, userResult := range userResults {
		user := userResult.ToModel()
		users = append(users, *user)
	}

	return users, nil
}

func (r *Repo) GetUserByID(ctx *context.Context, userId uuid.UUID) (*model.User, *model.ErrorResponse) {
	cCtx, cancel := context.WithTimeout(*ctx, r.queryTimeoutLimit)
	defer cancel()

	var userResult db.User
	args := map[string]any{
		"userId": userId,
	}

	err := r.queries.Read.getUser.GetContext(cCtx, &userResult, args)
	if err != nil {
		return nil, &model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	user := userResult.ToModel()

	return user, nil
}

func (r *Repo) Close() (err error) {
	err = r.queries.Close()
	if err != nil {
		return fmt.Errorf("error closing queries %w", err)
	}
	return
}
