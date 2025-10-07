package model

import (
	"time"

	"github.com/google/uuid"
)

type ProjectSponsor struct {
	ProjectID uuid.UUID
	SponsorID uuid.UUID
	StartDate *time.Time
	EndDate *time.Time
}