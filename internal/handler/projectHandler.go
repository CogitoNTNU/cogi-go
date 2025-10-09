package handler

import (
	"context"
	"net/http"

	"github.com/CogitoNTNU/cogi-go/internal/service"
	"github.com/gin-gonic/gin"
)

type Project struct {
	service *service.Project
	ctx     *context.Context
}

func NewProjectHandler(projectService *service.Project, ctx *context.Context) *Project {
	return &Project{service: projectService, ctx: ctx}
}

func (p *Project) GetAllProjects(gCtx *gin.Context) {
	allProjects, err := p.service.GetAllProjects(p.ctx)
	if err != nil {
		gCtx.AbortWithStatusJSON(err.Code, err.Message)
		return
	}

	gCtx.JSON(http.StatusOK, allProjects)
}

func (p *Project) GetProjectByID(gCtx *gin.Context) {
	projectId := gCtx.Param("projectId")
	project, err := p.service.GetProjectByID(p.ctx, projectId)

	if err != nil {
		gCtx.AbortWithStatusJSON(err.Code, err.Message)
		return
	}

	gCtx.JSON(http.StatusOK, project)
}
