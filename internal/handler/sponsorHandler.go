package handler

import (
	"context"
	"net/http"

	"github.com/CogitoNTNU/cogi-go/internal/service"
	"github.com/gin-gonic/gin"
)

type Sponsor struct {
	service *service.Sponsor
	ctx     *context.Context
}

func NewSponsorHandler(sponsorService *service.Sponsor, ctx *context.Context) *Sponsor {
	return &Sponsor{service: sponsorService, ctx: ctx}
}

func (s *Sponsor) GetAllSponsors(gCtx *gin.Context) {
	allSponsors, err := s.service.GetAllSponsors(s.ctx)
	if err != nil {
		gCtx.AbortWithStatusJSON(err.Code, err.Message)
		return
	}

	gCtx.JSON(http.StatusOK, allSponsors)
}

func (s *Sponsor) GetSponsorByID(gCtx *gin.Context) {
	sponsorId := gCtx.Param("id")
	sponsor, err := s.service.GetSponsorByID(s.ctx, sponsorId)
	if err != nil {
		gCtx.AbortWithStatusJSON(err.Code, err.Message)
		return
	}

	gCtx.JSON(http.StatusOK, sponsor)
}
