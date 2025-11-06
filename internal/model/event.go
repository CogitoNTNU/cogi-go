package model

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	WorkshopEventType EventType = "workshop"
	NewsEventType EventType = "news"
	HackathonEventType EventType = "hackathon"
	MeetingEventType EventType = "meeting"
	OtherEventType EventType = "other"


)

type Event struct {
	Id uuid.UUID
	Name string
	StartAt time.Time
	EndAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Type EventType
	Location string
	Description string
	Content []byte
	MaxAttendees int64
	Attendees []User
}
