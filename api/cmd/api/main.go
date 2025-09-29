package main

import (
	"fmt"
	"log"
	"os"

	"github.com/grundigdev/club/handlers"
	"github.com/grundigdev/club/mailer"
	"github.com/grundigdev/club/shared"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Application struct {
	logger  echo.Logger
	server  *echo.Echo
	handler handlers.Handler
}

func main() {
	e := echo.New()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := shared.NewPostgres()
	if err != nil {
		e.Logger.Fatal(err.Error())
	}

	appMailer := mailer.NewMailer(e.Logger)

	h := handlers.Handler{
		DB:     db,
		Logger: e.Logger,
		Mailer: appMailer,
	}
	app := Application{
		logger:  e.Logger,
		server:  e,
		handler: h,
	}
	e.Use(middleware.Logger())
	//e.Use(middlewares.CustomMiddleware)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://127.0.0.1:5174", "http://localhost:5174", "http://127.0.0.1:5173", "http://localhost:5173", "http://192.168.2.106:5173", "https://tiktok.grundigdev.de", "https://tiktok.grundigdev.de:443"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true, // If you're dealing with cookies
	}))
	app.routes(h)
	fmt.Println(app)
	port := os.Getenv("API_PORT")
	appAddress := fmt.Sprintf("localhost:%s", port)
	e.Logger.Fatal(e.Start(appAddress))

	//e.Logger.Fatal(e.Start("192.168.2.106:8000"))
}
