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

// NewRepo initializes a new user repository.
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

func (r *Repo) GetUserByID(ctx *context.Context, userID uuid.UUID) (*model.User, *model.ErrorResponse) {
	cCtx, cancel := context.WithTimeout(*ctx, r.queryTimeoutLimit)
	defer cancel()

	var userResult db.User
	args := map[string]any{
		"userId": userID,
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

func (r *Repo) Close() error {
	if err := r.queries.Close(); err != nil {
		return fmt.Errorf("error closing queries: %w", err)
	}
	return nil
}

// getUserByEmail retrieves a user by email address.
func (r *Repo) getUserByEmail(ctx *context.Context, email string) (*model.User, *model.ErrorResponse) {
	c, cancel := context.WithTimeout(*ctx, r.queryTimeoutLimit)
	defer cancel()

	var row struct {
		UserID          string     `db:"user_id"`
		FirstName       string     `db:"first_name"`
		LastName        string     `db:"last_name"`
		Description     string     `db:"description"`
		Nickname        string     `db:"nickname"`
		Email           string     `db:"email"`
		Phone           string     `db:"phone"`
		Gender          string     `db:"gender"`
		GithubURL       string     `db:"github_url"`
		LinkedinURL     string     `db:"linkedin_url"`
		KaggleURL       string     `db:"kaggle_url"`
		HuggingfaceURL  string     `db:"huggingface_url"`
		Password        string     `db:"password"`
		Avatar          string     `db:"avatar"`
		ImagePermission *time.Time `db:"image_permission"`
		FoodPreference  []string   `db:"food_preference"`
	}

	if err := r.queries.Read.getUserByEmail.GetContext(c, &row, map[string]any{"email": email}); err != nil {
		r.log.WithError(err).Warn("user not found by email")
		return nil, &model.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "user not found",
		}
	}

	user := &model.User{
		Id:              row.UserID,
		FirstName:       row.FirstName,
		LastName:        row.LastName,
		Description:     row.Description,
		Nickname:        row.Nickname,
		Email:           row.Email,
		Phone:           row.Phone,
		Gender:          row.Gender,
		GithubURL:       row.GithubURL,
		LinkedinURL:     row.LinkedinURL,
		KaggleURL:       row.KaggleURL,
		HuggingfaceURL:  row.HuggingfaceURL,
		Password:        row.Password,
		Avatar:          row.Avatar,
		ImagePermission: row.ImagePermission,
		FoodPreference:  row.FoodPreference,
	}

	return user, nil
}
