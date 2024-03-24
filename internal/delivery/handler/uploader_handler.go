package handler

import (
	"errors"
	"net/http"

	"github.com/Dzikuri/openidea-segokuning/internal/helper"
	"github.com/Dzikuri/openidea-segokuning/internal/model"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

func (h *Handler) ImageUpload(c echo.Context) error {

	file, err := c.FormFile("file")

	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	imageUrl, err := helper.UploadFileToS3(c.Request().Context(), uuid.NewV4().String(), file)
	if condition := err != nil; condition {

		if errors.Is(err, model.ErrFileSizeNotValid) {
			return c.JSON(http.StatusBadRequest, model.ResponseError{
				Code:    http.StatusBadRequest,
				Message: model.ErrFileSizeNotValid.Error(),
			})
		}

		if condition := errors.Is(err, model.ErrExtensionNotValid); condition {
			return c.JSON(http.StatusBadRequest, model.ResponseError{
				Code:    http.StatusBadRequest,
				Message: model.ErrExtensionNotValid.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, model.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, model.Response[any]{
		Data: map[string]interface{}{
			"imageUrl": imageUrl,
		},
		Message: "File uploaded successfully",
	})
}
