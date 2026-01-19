package model

import (
	"github.com/google/uuid"
)

type TeamMembers struct {
	TeamID uuid.UUID
	UserID uuid.UUID
	RoleID uuid.UUID
}
