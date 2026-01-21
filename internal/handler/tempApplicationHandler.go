package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/CogitoNTNU/cogi-go/internal/api/dto"
	"github.com/CogitoNTNU/cogi-go/internal/service"
	"github.com/CogitoNTNU/cogi-go/internal/util/env"
	"github.com/gin-gonic/gin"

	mail "github.com/xhit/go-simple-mail/v2"
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

	smtpClient := gCtx.MustGet("smtp_client").(*mail.SMTPClient)
	email := mail.NewMSG()
	email.SetFrom("From noreply@cogito-ntnu.no <noreply@cogito-ntnu.no>").AddTo(req.Email).SetSubject("Thank you for your project application to Cogito NTNU!")

	body := fmt.Sprintf("Dear %s,\n\nThank you for your application to join Cogito NTNU!\n We have registered your wish to join the projects: %s. Please send us a mail at styret@cogito-ntnu.no if something is wrong.\n\nBest Regards,\nCogito NTNU Board", req.FirstName, strings.Join(req.Projects, ", "))
	email.SetBody(mail.TextHTML, body)

	if email.Error != nil {
		gCtx.JSON(http.StatusInternalServerError, "Something went wrong when creating the reply email.")
		return
	}

	if email.Send(smtpClient) != nil {
		gCtx.JSON(http.StatusInternalServerError, "Something went wrong when sending the reply email.")
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
