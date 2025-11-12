package model

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	Id              uuid.UUID
	UserId          uuid.UUID
	ProjectId1      uuid.UUID
	ProjectId2      uuid.UUID
	ProjectId3      uuid.UUID
	ApplicationText string
	AppliedAt       time.Time
	ModifiedAt      time.Time
	IsMember        bool
}
