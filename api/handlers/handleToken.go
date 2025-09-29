package handlers

import (
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

	tokenService := services.NewTokenService(h.DB)

	token, err := tokenService.CreateToken(payload)
	if err != nil {
		return shared.SendInternalServerErrorResponse(c, err.Error())
	}

	return shared.SendSuccessResponse(c, "Token created", token)

}

func (h *Handler) CheckToken(c echo.Context) error {
	payload := new(requests.CheckTokenRequest)

	err := h.BindBodyRequest(c, payload)
	if err != nil {
		c.Logger().Error(err)
		return shared.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)

	if validationErrors != nil {
		return shared.SendFailedValidationResponse(c, validationErrors)
	}

	tokenService := services.NewTokenService(h.DB)
	token, isValid, err := tokenService.GetToken(payload.AccessToken)
	if err != nil {
		return shared.SendInternalServerErrorResponse(c, err.Error())
	}

	responseData := map[string]interface{}{
		"token":   token,
		"isValid": isValid,
	}

	return shared.SendSuccessResponse(c, "Token successfully fetched", responseData)
}
