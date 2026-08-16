package handler

import (
	"context"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/CogitoNTNU/cogi-go/internal/api/dto"
	"github.com/CogitoNTNU/cogi-go/internal/service"
	"github.com/CogitoNTNU/cogi-go/internal/util/env"
	"github.com/gin-gonic/gin"

	mail "github.com/xhit/go-simple-mail/v2"
)

const defaultApplicationDeadline = "2026-08-28T23:59:59+02:00"

type TempApplication struct {
	service *service.TempApplication
	ctx     *context.Context
	env     *env.EnvConfig
}

func NewTempApplicationHandler(tempApplicationService *service.TempApplication, ctx *context.Context, env *env.EnvConfig) *TempApplication {
	return &TempApplication{service: tempApplicationService, ctx: ctx, env: env}
}

func (t *TempApplication) CreateTempApplication(gCtx *gin.Context) {
	// Always enforce a deadline. Production may override the compiled Fall 2026
	// value, but a missing environment variable must never leave intake open.
	deadlineStr := defaultApplicationDeadline
	if configuredDeadline, err := t.env.Read("APPLICATION_DEADLINE"); err == nil && configuredDeadline != "" {
		deadlineStr = configuredDeadline
	}

	deadline, parseErr := time.Parse(time.RFC3339, deadlineStr)
	if parseErr != nil {
		fmt.Printf("Invalid APPLICATION_DEADLINE %q: %v\n", deadlineStr, parseErr)
		gCtx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "application intake is not configured"})
		return
	}
	if time.Now().After(deadline) {
		gCtx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "applications are closed"})
		return
	}

	// A real application is a few KB — cap the body so the public endpoint
	// can't be used to write multi-MB rows or relay huge emails.
	gCtx.Request.Body = http.MaxBytesReader(gCtx.Writer, gCtx.Request.Body, 64<<10)

	var req dto.CreateTempApplicationRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		// Missing return here used to let invalid/empty requests fall through
		// and insert empty rows in the database. Log rejections so an
		// application is never lost without a trace on either side.
		fmt.Printf("Rejected application payload from %s: %v\n", gCtx.ClientIP(), err)
		gCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := t.service.CreateTempApplication(t.ctx, &req)
	if err != nil {
		// err.Message can contain raw Postgres error text (schema/column names)
		// — log it, but never forward internals to an unauthenticated caller.
		fmt.Printf("Failed to store application from %s: %v\n", gCtx.ClientIP(), err.Message)
		gCtx.AbortWithStatusJSON(err.Code, gin.H{"error": "could not store the application"})
		return
	}

	// The application is stored at this point — email failures must not fail
	// the request, but they should be logged with the actual error.
	smtpServer := gCtx.MustGet("smtp_server").(*mail.SMTPServer)
	smtpClient, smtpErr := smtpServer.Connect()
	if smtpErr != nil {
		fmt.Printf("Could not connect to SMTP server: %v\n", smtpErr)
	} else {
		email := applicationReplyEmail(&req)
		if sendErr := email.Send(smtpClient); sendErr != nil {
			fmt.Printf("Could not send email to %s: %v\n", req.Email, sendErr)
		}
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
		fmt.Printf("CSV export failed: %v\n", errResp.Message)
		gCtx.AbortWithStatusJSON(errResp.Code, gin.H{"error": "export failed"})
		return
	}
}

func applicationReplyEmail(req *dto.CreateTempApplicationRequest) *mail.Email {
	email := mail.NewMSG()
	email.SetFrom("no-reply@cogito-ntnu.no").AddTo(req.Email).SetSubject("Thank you for your project application to Cogito NTNU!")

	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; line-height: 1.8; color: #333; max-width: 600px; margin: 0 auto; padding: 20px; }
        .logo { max-width: 150px; height: auto; }
        .content { background-color: #f9f9f9; padding: 25px; border-radius: 8px; }
				.projects { background-color: #1e90ff; padding: 15px; border-radius: 5px; margin: 20px 0; color: #ffffff; }
        .footer { text-align: center; margin-top: 30px; font-size: 12px; color: #888; }
    </style>
</head>
<body>
    <div class="content">
        <p>Dear %s,</p>
        <strong>Thank you for applying to Cogito NTNU!</strong>
        <p>We've received your application and we will review it shortly after the application deadline passes.</p>
        <div class="projects">
            <strong>The projects you applied to:</strong><br>
            %s<br><br> 
            <strong>Application text</strong><br>
            %s
        </div>
		<p>Best regards,<br><strong>The Cogito NTNU Board</strong></p>
		<p>Something not quite right? Send in a new application 😄</p>
    </div>
    <div class="footer">
		<p>This is an automatic reply. If you want to contact us, please use:</p>
        <p>styret@cogito-ntnu.no</p>
    </div>
</body>
</html>`, escapeForEmail(req.FirstName), strings.Join(escapeAllForEmail(req.Projects), "<br>"), escapeForEmail(req.ApplicationText))
	email.SetBody(mail.TextHTML, body)

	return email
}

// Applicant-supplied text goes into an HTML email sent from Cogito's official
// address to an applicant-chosen recipient — without escaping, anyone could
// make no-reply@cogito-ntnu.no deliver arbitrary HTML to any inbox.
func escapeForEmail(value string) string {
	return strings.ReplaceAll(html.EscapeString(value), "\n", "<br>")
}

func escapeAllForEmail(values []string) []string {
	escaped := make([]string, len(values))
	for i, value := range values {
		escaped[i] = escapeForEmail(value)
	}
	return escaped
}
