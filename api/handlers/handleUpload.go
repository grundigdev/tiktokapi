package handlers

import (
	"github.com/grundigdev/club/requests"
	"github.com/grundigdev/club/services"
	"github.com/grundigdev/club/shared"
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateUpload(c echo.Context) error {
	payload := new(requests.CreateUploadRequest)

	err := h.BindBodyRequest(c, payload)
	if err != nil {
		c.Logger().Error(err)
		return shared.SendBadRequestResponse(c, err.Error())
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)

	if validationErrors != nil {
		return shared.SendFailedValidationResponse(c, validationErrors)
	}

	uploadService := services.NewUploadService(h.DB)

	upload, err := uploadService.CreateUpload(payload)
	if err != nil {
		return shared.SendInternalServerErrorResponse(c, err.Error())
	}

	return shared.SendSuccessResponse(c, "Upload saved", upload)

}

func (h *Handler) GetUploads(c echo.Context) error {
	uploadService := services.NewUploadService(h.DB)
	uploads, err := uploadService.GetAllUploads()
	if err != nil {
		return shared.SendInternalServerErrorResponse(c, err.Error())
	}
	return shared.SendSuccessResponse(c, "All Uploads successfully fetched", uploads)
}
