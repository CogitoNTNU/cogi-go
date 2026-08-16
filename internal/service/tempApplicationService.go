package service

import (
	"context"
	"encoding/csv"
	"net/http"
	"strings"

	"github.com/CogitoNTNU/cogi-go/internal/api/dto"
	"github.com/CogitoNTNU/cogi-go/internal/model"
	tempApplicationRepository "github.com/CogitoNTNU/cogi-go/internal/repository/tempApplication"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TempApplication struct {
	logger     *logrus.Entry
	repository *tempApplicationRepository.Repo
}

func NewTempApplicationService(tempApplicationRepository *tempApplicationRepository.Repo, logger *logrus.Entry) *TempApplication {
	return &TempApplication{
		repository: tempApplicationRepository,
		logger:     logger,
	}
}

func (t *TempApplication) CreateTempApplication(ctx *context.Context, tempApplication *dto.CreateTempApplicationRequest) *model.ErrorResponse {
	err := t.repository.InsertTempApplication(ctx, tempApplication)
	if err != nil {
		t.logger.Errorf("An error has occured when creating a temporary application.")
		return err
	}

	return nil
}

// Excel and Sheets execute cells starting with these characters as formulas.
// Applicants control most CSV fields, so a hostile application text like
// "=HYPERLINK(...)" would otherwise run on a board member's machine.
func sanitizeCSVCell(value string) string {
	if value == "" {
		return value
	}
	switch value[0] {
	case '=', '+', '-', '@', '\t', '\r':
		return "'" + value
	}
	return value
}

func (t *TempApplication) ExportTempApplicationsCSV(ctx *context.Context, responseWriter gin.ResponseWriter) *model.ErrorResponse {
	tempApplications, errResp := t.repository.GetAllTempApplications(ctx)
	if errResp != nil {
		return errResp
	}

	csvData := make([][]string, len(tempApplications))
	for i, tempApplication := range tempApplications {
		csvData[i] = []string{
			sanitizeCSVCell(tempApplication.FirstName),
			sanitizeCSVCell(tempApplication.LastName),
			sanitizeCSVCell(tempApplication.Email),
			sanitizeCSVCell(tempApplication.PhoneNumber),
			sanitizeCSVCell(strings.Join(tempApplication.Projects, ", ")),
			sanitizeCSVCell(tempApplication.ApplicationText),
			tempApplication.CreatedAt.Local().String(),
		}
	}

	writer := csv.NewWriter(responseWriter)
	defer writer.Flush()

	header := []string{"FirstName", "LastName", "Email", "PhoneNumber", "Projects", "ApplicationText", "CreatedAt"}
	if err := writer.Write(header); err != nil {
		return &model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	for _, record := range csvData {
		if err := writer.Write(record); err != nil {
			return &model.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
	}

	if err := writer.Error(); err != nil {
		return &model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}
