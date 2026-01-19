package model

import (
	"time"

	"github.com/google/uuid"
)

type GenderType string

const (
	GenderTypeFemale    GenderType = "female"
	GenderTypeMale      GenderType = "male"
	GenderTypeNonBinary GenderType = "non_binary"
	GenderTypeOther     GenderType = "other"
)

type User struct {
	Id              uuid.UUID
	Nickname        string
	Description     string
	Password        string
	FirstName       string
	LastName        string
	Email           string
	Phone           string
	Gender          GenderType
	GithubUrl       *string
	LinkedinUrl     *string
	KaggleUrl       *string
	HuggingfaceUrl  *string
	Avatar          *string
	ImagePermission *time.Time
	FoodPreference  []string
}
