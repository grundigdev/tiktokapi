package handlers

import (
	"errors"

	"github.com/grundigdev/club/mailer"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	DB     *gorm.DB
	Logger echo.Logger
	Mailer mailer.Mailer
}

func (h *Handler) BindBodyRequest(c echo.Context, payload interface{}) error {
	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err != nil {
		c.Logger().Error(err)
		return errors.New("Failed to bind Body, make sure you are sending a valid payload " + err.Error())
	}
	return nil

}
