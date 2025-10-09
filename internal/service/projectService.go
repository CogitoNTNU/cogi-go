package service

import (
	"context"
	"net/http"

	"github.com/CogitoNTNU/cogi-go/internal/model"
	projectRepository "github.com/CogitoNTNU/cogi-go/internal/repository/project"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Project struct {
	logger *logrus.Entry
	repository *projectRepository.Repo
}

func NewProjectService(projectRepository *projectRepository.Repo) *Project {
	return &Project{repository: projectRepository}
}


func (p *Project) GetAllProjects(ctx *context.Context) ([]model.Project, *model.ErrorResponse) {
	allProjects, err := p.repository.GetAllProjects(ctx)
	if err != nil {
		p.logger.Errorf("An error has occured when retrieving all projects. Error code: %s", err.Code)
		return nil, err
	}

	return allProjects, nil
}

func (p *Project) GetProjectByID(ctx *context.Context, projectId string) (*model.Project, *model.ErrorResponse) {
	id, err := uuid.Parse(projectId)

	if err != nil {
		return nil, &model.ErrorResponse{
			Code: http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	project, fetchErr := p.repository.GetProjectByID(ctx, id)
	if fetchErr != nil {
		return nil, fetchErr
	}

	return project, nil
}
