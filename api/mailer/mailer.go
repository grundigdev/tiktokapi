package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gopkg.in/gomail.v2"
)

//go:embed templates/*
var templateFS embed.FS

type Mailer struct {
	dialer *gomail.Dialer
	sender string // no-reply@app.com
	Logger echo.Logger
}

type EmailData struct {
	AppName string
	Subject string
	Meta    interface{}
}

func NewMailer(logger echo.Logger) Mailer {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mailPort, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		logger.Fatal(err)
	}
	mailHost := os.Getenv("MAIL_HOST")
	mailUsername := os.Getenv("MAIL_USERNAME")
	mailPassword := os.Getenv("MAIL_PASSWORD")

	dialer := gomail.NewDialer(mailHost, mailPort, mailUsername, mailPassword)
	return Mailer{
		dialer: dialer,
		sender: mailUsername,
		Logger: logger,
	}
}

func (mailer *Mailer) Send(recipient string, templateFile string, data EmailData) error {
	absolutePath := fmt.Sprintf("templates/%s", templateFile)
	tmpl, err := template.ParseFS(templateFS, absolutePath)
	fmt.Println(absolutePath)
	if err != nil {
		mailer.Logger.Error(err)
		return err
	}
	data.AppName = os.Getenv("APP_NAME")
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		mailer.Logger.Error(err)
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		mailer.Logger.Error(err)
		return err
	}

	goMailMessage := gomail.NewMessage()
	goMailMessage.SetHeader("To", recipient)
	goMailMessage.SetHeader("From", mailer.sender)
	goMailMessage.SetHeader("Subject", subject.String())

	goMailMessage.SetBody("text/html", htmlBody.String())

	err = mailer.dialer.DialAndSend(goMailMessage)

	if err != nil {
		mailer.Logger.Error(err)
	}
	return nil
}
