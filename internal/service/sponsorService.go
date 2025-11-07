package service

import (
	"context"
	"net/http"

	"github.com/CogitoNTNU/cogi-go/internal/model"
	sponsorRepository "github.com/CogitoNTNU/cogi-go/internal/repository/sponsor"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Sponsor struct {
	logger *logrus.Entry
	repository *sponsorRepository.Repo
}

func NewSponsorService(sponsorRepository *sponsorRepository.Repo) *Sponsor {
	return &Sponsor{repository: sponsorRepository}
}

func (s *Sponsor) GetAllSponsors(ctx *context.Context) ([]model.Sponsor, *model.ErrorResponse) {
	allSponsors, err := s.repository.GetAllSponsors(ctx)
	if err != nil {
		s.logger.Errorf("An error has occured when retrieving all sponsors. Error code: %s", err.Code)
		return nil, err
	}

	return allSponsors, nil
}

func (s *Sponsor) GetSponsorByID(ctx *context.Context, sponsorId string) (*model.Sponsor, *model.ErrorResponse) {
	id, err := uuid.Parse(sponsorId)

	if err != nil {
		return nil, &model.ErrorResponse{
			Code: http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	sponsor, fetchErr := s.repository.GetSponsorByID(ctx, id)
	if fetchErr != nil {
		return nil, fetchErr
	}

	return sponsor, nil
}