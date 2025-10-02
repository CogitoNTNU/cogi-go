package db

import (
	"time"

	"github.com/CogitoNTNU/cogi-go/internal/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID              uuid.UUID     `db:"user_id"`
	FirstName       string         `db:"first_name"`
	LastName        string         `db:"last_name"`
	Description     string        `db:"description"`
	Nickname        string         `db:"nickname"`
	Email           string         `db:"email"`
	Phone           string         `db:"phone"`
	Gender          model.GenderType     `db:"gender"` 
	GithubURL       *string        `db:"github_url"`
	LinkedinURL     *string        `db:"linkedin_url"`
	KaggleURL       *string        `db:"kaggle_url"`
	HuggingfaceURL  *string        `db:"huggingface_url"`
	Password       string         `db:"password"`
	Avatar          *string        `db:"avatar"`
	ImagePermission *time.Time     `db:"image_permission"` 
	FoodPreference  pq.StringArray `db:"food_preference"` 
}

func (u *User) ToModel() *model.User {
	return &model.User{
		Id:             u.ID,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Description:    u.Description,
		Nickname:       u.Nickname,
		Email:          u.Email,
		Phone:          u.Phone,
		Gender:        u.Gender,
		GithubUrl:     u.GithubURL,
		LinkedinUrl:   u.LinkedinURL,
		KaggleUrl:     u.KaggleURL,
		HuggingfaceUrl: u.HuggingfaceURL,
		Avatar:        u.Avatar,
		ImagePermission: u.ImagePermission,
		FoodPreference: u.FoodPreference,
	}
}