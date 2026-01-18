package model

import (
	"time"
	"github.com/google/uuid"
)

type TempApplication struct {
	Id              uuid.UUID
	Email           string
	PhoneNumber     int
	Projects		[]string
	ApplicationText string
	CreatedAt      *time.Time
}