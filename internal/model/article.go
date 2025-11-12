package model

import (
	"time"

	"github.com/google/uuid"
)

type Article struct {
	Id        uuid.UUID
	Title     string
	Content   []byte
	AuthorId  uuid.UUID
	ProjectId uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}
