package routes

import (
	"github.com/CogitoNTNU/cogi-go/internal/handler"
	"github.com/gin-gonic/gin"
)

// TODO: Smart way to inject all handlers
func RegisterPublicRoutes(router *gin.Engine, handlers *handler.User) {
	// TODO: Initialize public routes here
	// e.g router.GET("/public", publicHandler)

	router.GET("/users", handlers.GetAllUsers)
	router.GET("/users/:userId", handlers.GetUserByID)

}

func RegisterPrivateRoutes(router *gin.Engine) {

}

func RegisterAdminRoutes(router *gin.Engine) {

}
