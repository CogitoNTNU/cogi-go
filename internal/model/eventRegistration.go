package model

import (
	"time"

	"github.com/google/uuid"
)

type EventRegistration struct {
	ID uuid.UUID
	EventId uuid.UUID
	UserId uuid.UUID

 	CreatedAt time.Time
}
