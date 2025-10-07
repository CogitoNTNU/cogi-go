package model

import (
	"github.com/google/uuid"
)

type SponsorLevel string

const (
	SponsorLevelBronze SponsorLevel = "bronze"
	SponsorLevelSilver  SponsorLevel = "silver"
	SponsorLevelGold    SponsorLevel = "gold"
)

type Sponsor struct {
	Id        uuid.UUID
	Name      string
	Logo      *string
	Website   *string
	Description string
	Level     SponsorLevel
}
