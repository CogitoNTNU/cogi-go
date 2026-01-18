package service

import (
	"context"
	"net/http"
	"fmt"
	"encoding/csv"
	"github.com/gin-gonic/gin"

	"github.com/CogitoNTNU/cogi-go/internal/model"
	tempApplicationRepository "github.com/CogitoNTNU/cogi-go/internal/repository/tempApplication"
	"github.com/sirupsen/logrus"
)

type TempApplication struct {
	logger     *logrus.Entry
	repository *tempApplicationRepository.Repo
}

func NewTempApplicationService(tempApplicationRepository *tempApplicationRepository.Repo) *TempApplication {
	return &TempApplication{repository: tempApplicationRepository}
}

func (ta *TempApplication) ExportTempApplicationsCSV(ctx *context.Context, responseWriter gin.ResponseWriter) (*model.ErrorResponse) {
	tempApplications, errResp := ta.repository.GetAllTempApplications(ctx)
	if errResp != nil {
		return errResp
	}

	csvData := make([][]string, len(tempApplications))
	for i, tempApplication := range tempApplications {
		csvData[i] = []string{
			tempApplication.Email,
			tempApplication.PhoneNumber,
			fmt.Sprintf("%v", tempApplication.Projects), // Convert slice to string
			tempApplication.ApplicationText,
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

func (ta *TempApplication) InsertTempApplication(ctx *context.Context, tempApplication *model.TempApplication) *model.ErrorResponse {
	errResp := ta.repository.InsertTempApplication(ctx, tempApplication)
	if errResp != nil {
		ta.logger.Errorf("An error has occured when inserting a temp application. Error code: %s", errResp.Code)
		return errResp
	}

	return nil
}

func (ta *TempApplication) GetAllTempApplicationsCSV(ctx *context.Context) ([]model.TempApplication,*model.ErrorResponse) {
	tempApplications, errResp := ta.repository.GetAllTempApplications(ctx)
	if errResp != nil {
		ta.logger.Errorf("An error has occured when exporting temp applications to CSV. Error code: %s", errResp.Code)
		return nil, errResp
	}

	return tempApplications, nil
}
