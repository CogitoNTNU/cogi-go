package model

import (
	"time"

	"github.com/google/uuid"
)

type RegistrationStatus string

const (
    RegistrationRequested RegistrationStatus = "requested"
    RegistrationConfirmed RegistrationStatus = "confirmed"
    RegistrationWaitlisted RegistrationStatus = "waitlisted"
    RegistrationCancelled RegistrationStatus = "cancelled"
)

type EventRegistration struct {
	ID uuid.UUID
	EventId uuid.UUID
	UserId uuid.UUID

	Status RegistrationStatus

 	CreatedAt time.Time
    UpdatedAt time.Time

    Event *Event
    User *User

}
