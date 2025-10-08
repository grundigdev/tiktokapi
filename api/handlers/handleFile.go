package handlers

import (
	"github.com/grundigdev/club/requests"
	"github.com/grundigdev/club/services"
	"github.com/grundigdev/club/shared"
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateFile(c echo.Context) error {
	payload := new(requests.CreateFileRequest)

	err := h.BindBodyRequest(c, payload)
	if err != nil {
		c.Logger().Error(err)
		return shared.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)

	if validationErrors != nil {
		return shared.SendFailedValidationResponse(c, validationErrors)
	}

	fileService := services.NewFileService(h.DB)

	token, err := fileService.CreateFile(payload)
	if err != nil {
		return shared.SendInternalServerErrorResponse(c, err.Error())
	}

	return shared.SendSuccessResponse(c, "File created", token)
}

func (h *Handler) GetFile(c echo.Context) error {
	payload := new(requests.GetFileRequest)

	err := h.BindBodyRequest(c, payload)
	if err != nil {
		c.Logger().Error(err)
		return shared.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)

	if validationErrors != nil {
		return shared.SendFailedValidationResponse(c, validationErrors)
	}

	fileService := services.NewFileService(h.DB)
	file, err := fileService.GetFile(payload.ID)
	if err != nil {
		return shared.SendInternalServerErrorResponse(c, err.Error())
	}

	return shared.SendSuccessResponse(c, "File successfully fetched", file)
}

func (h *Handler) UpdateFile(c echo.Context) error {
	payload := new(requests.CreateFileRequest)

	err := h.BindBodyRequest(c, payload)
	if err != nil {
		c.Logger().Error(err)
		return shared.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)

	if validationErrors != nil {
		return shared.SendFailedValidationResponse(c, validationErrors)
	}

	fileService := services.NewFileService(h.DB)
	file, err := fileService.UpdateFile(payload)
	if err != nil {
		return shared.SendInternalServerErrorResponse(c, err.Error())
	}

	return shared.SendSuccessResponse(c, "File successfully changed", file)
}
