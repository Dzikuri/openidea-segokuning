package handler

import (
	"errors"
	"net/http"

	"github.com/Dzikuri/openidea-segokuning/internal/model"
	"github.com/labstack/echo/v4"
)

func (h *Handler) UserRegister(c echo.Context) error {
	var request model.UserAuthRequest
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
			Code:    model.ErrResBadRequest.Code,
			Message: model.ErrResBadRequest.Message,
			Error:   err,
		})
	}

	err = request.Validate()
	if err != nil {
		return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
			Code:    model.ErrResBadRequest.Code,
			Message: model.ErrResBadRequest.Message,
			Error:   err,
		})
	}

	result, err := h.UseCase.UserRegister(c.Request().Context(), &request)
	if err != nil {

		if errors.Is(err, model.ErrUserAlreadyExists) {
			return c.JSON(echo.ErrConflict.Code, model.ResponseError{
				Code:    model.ErrResBadRequest.Code,
				Message: model.ErrUserAlreadyExists.Error(),
				Error:   err,
			})
		}

		if errors.Is(err, echo.ErrInternalServerError) {
			return c.JSON(echo.ErrInternalServerError.Code, model.ResponseError{
				Code:    echo.ErrInternalServerError.Code,
				Message: echo.ErrInternalServerError.Error(),
				Error:   err,
			})
		}

		return c.JSON(echo.ErrInternalServerError.Code, model.ResponseError{
			Code:    echo.ErrInternalServerError.Code,
			Message: err.Error(),
			Error:   err,
		})
	}

	return c.JSON(http.StatusCreated, model.Response[any]{
		Code:    http.StatusCreated,
		Data:    result,
		Message: "Success",
	})
}

func (h *Handler) UserLogin(c echo.Context) error {
	var request model.UserLoginRequest
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
			Code:    model.ErrResBadRequest.Code,
			Message: model.ErrResBadRequest.Message,
			Error:   err,
		})
	}

	err = request.Validate()
	if err != nil {
		return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
			Code:    model.ErrResBadRequest.Code,
			Message: model.ErrResBadRequest.Message,
			Error:   err,
		})
	}

	result, err := h.UseCase.UserLogin(c.Request().Context(), &request)

	if err != nil {

		if errors.Is(err, model.ErrPasswordNotMatch) {
			return c.JSON(echo.ErrBadRequest.Code, model.ResponseError{
				Code:    echo.ErrBadRequest.Code,
				Message: model.ErrPasswordNotMatch.Error(),
				Error:   err,
			})
		}

		if errors.Is(err, model.ErrUserNotFound) {
			return c.JSON(echo.ErrNotFound.Code, model.ResponseError{
				Code:    echo.ErrNotFound.Code,
				Message: model.ErrUserNotFound.Error(),
				Error:   err,
			})
		}

		return c.JSON(echo.ErrInternalServerError.Code, model.ResponseError{
			Code:    echo.ErrInternalServerError.Code,
			Message: err.Error(),
			Error:   err,
		})
	}

	return c.JSON(http.StatusOK, model.Response[any]{
		Code:    http.StatusOK,
		Data:    result,
		Message: "User logged successfully",
	})

}
