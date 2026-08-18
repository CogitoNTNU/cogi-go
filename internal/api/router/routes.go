package routes

import (
	"github.com/CogitoNTNU/cogi-go/internal/handler"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	User            *handler.User
	Project         *handler.Project
	Sponsor         *handler.Sponsor
	TempApplication *handler.TempApplication
	S3Tester        *handler.S3Tester
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

	sponsors := api.Group("/sponsors")
	{
		sponsors.GET("", handlers.Sponsor.GetAllSponsors)
		sponsors.GET("/:id", handlers.Sponsor.GetSponsorByID)
	}

	tempApplication := api.Group("/temp-member-application")
	{
		tempApplication.POST("/", handlers.TempApplication.CreateTempApplication)
		tempApplication.GET("/export-csv", handlers.TempApplication.ExportTempApplicationsCSV)
	}

	s3Tester := api.Group("/s3-tester")
	{
		s3Tester.POST("/upload", handlers.S3Tester.TestUpload)
		s3Tester.POST("/upload-large", handlers.S3Tester.TestUploadLarge)
		s3Tester.POST("/delete", handlers.S3Tester.TestDelete)
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
