package main

import (
	"github.com/grundigdev/club/handlers"
)

func (app *Application) routes(handler handlers.Handler) {
	app.server.GET("/", handler.HealthCheck)

	apiGroup := app.server.Group("/api")

	absenceRoutes := apiGroup.Group("/token")
	{
		absenceRoutes.POST("/create", handler.CreateToken)
		absenceRoutes.POST("/get", handler.CheckToken)
	}

}
