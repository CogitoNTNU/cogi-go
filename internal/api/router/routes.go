package routes

import (
	"github.com/CogitoNTNU/cogi-go/internal/handler"
	jwt "github.com/appleboy/gin-jwt/v3"
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

func RegisterPrivateRoutes(router *gin.Engine, register *jwt.GinJWTMiddleware, handlers *Handlers) {
	handlerfunc := register.MiddlewareFunc()
	r := router.Group("/api", handlerfunc)
	auth := r.Group("/auth")
	{
		auth.GET("/profile", handlers.User.GetUserByID)
	}
}

func RegisterAdminRoutes(router *gin.Engine) {

}
