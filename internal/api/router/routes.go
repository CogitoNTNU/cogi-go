package routes

import (
	"github.com/CogitoNTNU/cogi-go/internal/handler"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	User    *handler.User
	Project *handler.Project
	AuthMW *jwt.GinJWTMiddleware
	GoogleLogin gin.HandlerFunc
}

func RegisterPublicRoutes(router *gin.Engine, handlers *Handlers) {
	api := router.Group("/api")
	auth := api.Group("/auth") 
	{
		auth.POST("/login", handlers.AuthMW.LoginHandler)
		auth.POST("/logout", handlers.AuthMW.LogoutHandler)
		auth.GET("/refresh", handlers.AuthMW.RefreshHandler)
		auth.POST("/google", handlers.GoogleLogin)
	}

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
	handlerFunc := register.MiddlewareFunc()
	api := router.Group("/api", handlerFunc)
	r := router.Group("/api", handlerFunc)
	api.Use(handlers.AuthMW.MiddlewareFunc())
	auth := r.Group("/auth")
	{
		auth.GET("/profile", handlers.User.GetUserByID)
	}
}

func RegisterAdminRoutes(router *gin.Engine) {
	// TODO: Implement admin-specific routes here.
}