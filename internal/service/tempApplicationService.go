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

func (t *TempApplication) ExportTempApplicationsCSV(ctx *context.Context, responseWriter gin.ResponseWriter) *model.ErrorResponse {
	tempApplications, errResp := t.repository.GetAllTempApplications(ctx)
	if errResp != nil {
		return errResp
	}

	csvData := make([][]string, len(tempApplications))
	for i, tempApplication := range tempApplications {
		csvData[i] = []string{
			tempApplication.FirstName,
			tempApplication.LastName,
			tempApplication.Email,
			tempApplication.PhoneNumber,
			strings.Join(tempApplication.Projects, ", "),
			tempApplication.ApplicationText,
			tempApplication.CreatedAt.Local().String(),
		}
	}

	writer := csv.NewWriter(responseWriter)
	defer writer.Flush()

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
