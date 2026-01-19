package model

import (
	"time"

	"github.com/google/uuid"
)

type Achivement struct {
	Id      uuid.UUID
	Title   string
	IconUrl *string
}

type UserAchivement struct {
	Id           uuid.UUID
	UserId       uuid.UUID
	AchivementId uuid.UUID
	DateEarned   *time.Time
}
