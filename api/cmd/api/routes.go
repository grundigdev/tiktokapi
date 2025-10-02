package main

import (
	"github.com/grundigdev/club/handlers"
)

func (app *Application) routes(handler handlers.Handler) {
	app.server.GET("/health", handler.HealthCheck)

	apiGroup := app.server.Group("/api")

	tokenRoutes := apiGroup.Group("/token")
	{
		tokenRoutes.POST("/create", handler.CreateToken)
		tokenRoutes.POST("/get", handler.CheckToken)
	}

	uploadRoutes := apiGroup.Group("/upload")
	{
		uploadRoutes.POST("/create", handler.CreateUpload)
		uploadRoutes.POST("/get/all", handler.GetUploads)
	}

}
