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

func (r*Repo) getUserByEmail(ctx *context.Context, email string) (*model.User, *model.ErrorResponse) {
	c, cancel := context.WithTimeout(*ctx, r.timeout)
	defer cancel()
var row stuct {
	UserID string `db:"user_id"`
	FirstName string `db:"first_name"`
	LastName string `db:"last_name"`
	Description string `db:"description"`
	Nickname string `db:"nickname"`
	Email string `db:"email"`
	Phone string `db:"phone"`
	Gender string `db:"gender"`
	GithubURL string `db:"github_url"`
	LinkedinURL string `db:"linkedin_url"`
	KaggleURL string `db:"kaggle_url"`
	HuggingfaceURL string `db:"huggingface_url"`
	Password string `db:"password"`
	Avatar string `db:"avatar"`
	ImagePermission *time.Time `db:"image_permission"` //!
	FoodPrefrence []string `db:"food_prefrence"`
}
if err := r.queries.Read.getUserByEmail.GetContext(c, &row, map[string]any{"email": email}); err != nil {
	r.logger.WithError(err).Warn("user not found by email")
	return nil, &model.ErrorResponse{Code: http.StatusNotFound, Message: "user not found"}
}
user := &model.User{
	ID: row.UserID,
	FirstName: row.FirstName,
	LastName: row.LastName,
	Description: row.Description,
	Nickname: row.Nickname,
	Email: row.Email,
	Phone: row.Phone,
	Gender: row.Gender,
	GithubUrl: row.GithubURL,
	linkedinUrl: row.LinkedinURL,
	KaggleUrl: row.KaggleURL,
	HuggingfaceUrl: row.HuggingfaceURL,
	Password: row.Password,
	Avatar: row.Avatar
	ImagePermission: row.ImagePermission,
	FoodPrefrence: row.FoodPrefrence,
}
return user, nil
}