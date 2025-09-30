package routes

import "github.com/gin-gonic/gin"


// TODO: Smart way to inject all handlers
func RegisterPublicRoutes(router *gin.Engine, ) {
	// TODO: Initialize public routes here
	// e.g router.GET("/public", publicHandler)

	router.GET("/", )

}

func RegisterPrivateRoutes(router *gin.Engine) {

}

func RegisterAdminRoutes(router *gin.Engine) {

}
