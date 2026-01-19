package model

import (
	"github.com/google/uuid"
)

type Project struct {
	Id            uuid.UUID
	Title         string
	GithubUrl     *string
	Logo          *string
	Playable      bool
	Released      bool
	ActiveProject bool
	ProjectUrl    *string
}
