package model

import (
	"github.com/google/uuid"
)

type Semester string;
type GroupType string;

const (
	SpringSemester Semester = "spring"
	AutumnSemester Semester = "autumn"
)

const (
	MarketingGroup GroupType = "marketing"
	BoardGroup    GroupType = "board"
	SocialGroup   GroupType = "social"
)

type Team struct {
	Id          uuid.UUID
	Semester    Semester
	Year       int
}

type ProjectTeam struct {
	Team
	ProjectId uuid.UUID
}

type AdministrasjonTeam struct {
	Team
	GroupType GroupType
}