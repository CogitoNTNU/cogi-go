package routes

import (
	"github.com/CogitoNTNU/cogi-go/internal/handler"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	User    *handler.User
	Project *handler.Project
}

func RegisterPublicRoutes(router *gin.Engine, handlers *Handlers) {
	api := router.Group("/api")
    
    users := api.Group("/users")
    {
        users.GET("", handlers.User.GetAllUsers)
        users.GET("/:userId", handlers.User.GetUserByID)
    }
    
    projects := api.Group("/projects")
    {
        projects.GET("", handlers.Project.GetAllProjects)
        projects.GET("/:projectId", handlers.Project.GetProjectByID)
    }
}

func RegisterPrivateRoutes(router *gin.Engine) {

}

func RegisterAdminRoutes(router *gin.Engine) {

}
