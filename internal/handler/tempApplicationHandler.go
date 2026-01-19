package handler

import (
	"context"
	"net/http"

	"github.com/CogitoNTNU/cogi-go/internal/api/dto"
	"github.com/CogitoNTNU/cogi-go/internal/service"
	"github.com/CogitoNTNU/cogi-go/internal/util/env"
	"github.com/gin-gonic/gin"
)

type TempApplication struct {
	service *service.TempApplication
	ctx     *context.Context
	env     *env.EnvConfig
}

func NewTempApplicationHandler(tempApplicationService *service.TempApplication, ctx *context.Context, env *env.EnvConfig) *TempApplication {
	return &TempApplication{service: tempApplicationService, ctx: ctx, env: env}
}

func (t *TempApplication) CreateTempApplication(gCtx *gin.Context) {
	var req dto.CreateTempApplicationRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
	}

	err := t.service.CreateTempApplication(t.ctx, &req)
	if err != nil {
		gCtx.AbortWithStatusJSON(err.Code, err.Message)
		return
	}

	gCtx.Status(http.StatusCreated)
}

func (t *TempApplication) ExportTempApplicationsCSV(gCtx *gin.Context) {
	// WARNING: TEMPORARY AUTH CHECK
	// TODO: Refactor this with Gin Middleware
	token := gCtx.GetHeader("X-Admin-Token")
	expectedToken, err := t.env.Read("ADMIN_TOKEN")
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "configuration error"})
		return
	}

	if token != expectedToken || expectedToken == "" {
		gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	gCtx.Header("Content-Type", "text/csv")
	gCtx.Header("Content-Disposition", "attachment; filename=member_applications.csv")

	errResp := t.service.ExportTempApplicationsCSV(t.ctx, gCtx.Writer)
	if errResp != nil {
		gCtx.AbortWithStatusJSON(errResp.Code, errResp.Message)
		return
	}
}
