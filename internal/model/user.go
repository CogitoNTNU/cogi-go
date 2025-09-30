package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id             uuid.UUID
	Nickname       string
	Password       string
	FirstName      string
	LastName       string
	Email          string
	Phone          string
	GithubUrl      string
	LinkedinUrl    string
	Avatar         string
	PhotoPermit    time.Time
	FoodPreference []string
	IsActive       bool
}
