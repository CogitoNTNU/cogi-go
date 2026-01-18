package model

import (
	"time"
	"github.com/google/uuid"
)

type TempApplication struct {
	Id              uuid.UUID
	FirstName       string
	LastName        string
	Email           string
	PhoneNumber     string
	Projects        []string
	ApplicationText string
	CreatedAt      *time.Time
}