package handlers

import (
	"fmt"

	"github.com/grundigdev/club/requests"
	"github.com/grundigdev/club/services"
	"github.com/grundigdev/club/shared"
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateToken(c echo.Context) error {
	payload := new(requests.CreateTokenRequest)

	err := h.BindBodyRequest(c, payload)
	if err != nil {
		c.Logger().Error(err)
		return shared.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)

	if validationErrors != nil {
		return shared.SendFailedValidationResponse(c, validationErrors)
	}
	fmt.Println(payload)
	tokenService := services.NewTokenService(h.DB)

	token, err := tokenService.CreateToken(payload)
	if err != nil {
		return shared.SendInternalServerErrorResponse(c, err.Error())
	}

	return shared.SendSuccessResponse(c, "Token created", token)

}

func (h *Handler) GetLastToken(c echo.Context) error {
	tokenService := services.NewTokenService(h.DB)
	lastToken, err := tokenService.GetLastToken()
	if err != nil {
		return shared.SendInternalServerErrorResponse(c, err.Error())
	}
	return shared.SendSuccessResponse(c, "Last Token successfully fetched", lastToken)
}
