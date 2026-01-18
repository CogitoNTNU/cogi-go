package handler

import (
	"context"
	"net/http"
	"os"

	"github.com/CogitoNTNU/cogi-go/internal/service"
	"github.com/gin-gonic/gin"
)

type TempApplication struct {
	service *service.TempApplicationService
	ctx     *context.Context
}

func NewTempApplication(userService *service.User, ctx *context.Context) *User {
	return &User{service: userService, ctx: ctx}
}

func (t *TempApplication) CreateTempApplication(gCtx *gin.Context) {
	tempApplication, err := t.service.CreateTempApplication(t.ctx)
	if err != nil {
		gCtx.AbortWithStatusJSON(err.Code, err.Message)
		return
	}

	gCtx.JSON(http.StatusOK, tempApplication)
}

func (t *TempApplication) ExportTempApplicationsCSV(gCtx *gin.Context) {
	// WARNING: TEMPORARY AUTH CHECK
	// TODO: Refactor this with Gin Middleware
	token := gCtx.GetHeader("X-Admin-Token")
	expectedToken := os.Getenv("ADMIN_EXPORT_TOKEN")

	if token != expectedToken || expectedToken == "" {
		gCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	gCtx.Header("Content-Type", "text/csv")
	gCtx.Header("Content-Disposition", "attachment; filename=member_applications.csv")

	err := service.ExportTempApplicationsCSV(t.ctx, gCtx.Writer)
	if err != nil {
		gCtx.AbortWithStatusJSON(err.Code, err.Message)
		return
	}
}
