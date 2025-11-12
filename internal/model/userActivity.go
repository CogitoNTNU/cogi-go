package model

import (
	"time"

	"github.com/google/uuid"
)

type UserActivity struct {
	Id       uuid.UUID
	UserId   uuid.UUID
	Activity []time.Time
}
