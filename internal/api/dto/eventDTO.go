package dto

import (
	"time"

	"github.com/CogitoNTNU/cogi-go/internal/model"
)

type CreateEventRequest struct {
	Name         string          `json:"name"`
	StartAt      time.Time       `json:"startAt"`
	EndAt        time.Time       `json:"endAt"`
	Type         model.EventType `json:"type"`
	Location     string          `json:"location"`
	Description  string          `json:"description"`
	Content      []byte          `json:"content"`
	MaxAttendees int             `json:"maxAttendees"`
}
