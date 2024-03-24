package handler

import (
	"errors"
	"net/http"

	"github.com/Dzikuri/openidea-segokuning/internal/model"
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateFriend(c echo.Context) error {
	var request model.FriendRequest
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

	usr, ok := c.Get("userId").(*model.UserResponse)
	if ok {
		request.FriendId = request.UserId
		request.UserId = usr.Id.String()
	}

	if !ok {
		return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
			Code:    model.ErrResBadRequest.Code,
			Message: model.ErrResBadRequest.Message,
			Error:   err,
		})
	}

	if request.FriendId == request.UserId {
		return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
			Code:    model.ErrResBadRequest.Code,
			Message: "Cannot Add Yourself as Friend",
			Error:   err,
		})
	}

	_, err = h.UseCase.AddFriend(c.Request().Context(), request.UserId, request.FriendId)

	if err != nil {

		if condition := errors.Is(err, model.ErrInvalidUserId); condition {

			return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
				Code:    model.ErrResBadRequest.Code,
				Message: model.ErrInvalidUserId.Error(),
				Error:   err,
			})
		}

		if condition := errors.Is(err, model.ErrUserNotFound); condition {
			return c.JSON(model.ErrResNotFound.Code, model.ResponseError{
				Code:    model.ErrResNotFound.Code,
				Message: model.ErrUserNotFound.Error(),
				Error:   err,
			})

		}

		if condition := errors.Is(err, model.ErrResNotFound.Error); condition {
			return c.JSON(model.ErrResNotFound.Code, model.ResponseError{
				Code:    model.ErrResNotFound.Code,
				Message: model.ErrUserNotFound.Error(),
				Error:   err,
			})

		}

		if condition := errors.Is(err, model.ErrAlreadyBeFriend); condition {
			return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
				Code:    model.ErrResBadRequest.Code,
				Message: model.ErrAlreadyBeFriend.Error(),
				Error:   err,
			})
		}

		return c.JSON(echo.ErrInternalServerError.Code, model.ResponseError{
			Code:    echo.ErrInternalServerError.Code,
			Message: echo.ErrInternalServerError.Error(),
			Error:   err,
		})
	}

	return c.JSON(http.StatusOK, model.Response[any]{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    make(map[string]interface{}),
	})
}

func (h *Handler) GetFriends(c echo.Context) error {
	var request model.GetFriendListRequest

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

	usr, ok := c.Get("userId").(*model.UserResponse)
	if ok {
		request.UserId = usr.Id.String()
	}

	if !ok {
		return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
			Code:    model.ErrResBadRequest.Code,
			Message: model.ErrResBadRequest.Message,
			Error:   err,
		})
	}

	result, err := h.UseCase.GetFriendList(c.Request().Context(), request)

	if err != nil {
		return c.JSON(echo.ErrInternalServerError.Code, model.ResponseError{
			Code:    echo.ErrInternalServerError.Code,
			Message: echo.ErrInternalServerError.Error(),
			Error:   err,
		})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *Handler) DeleteFriend(c echo.Context) error {
	var request model.FriendRequest
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

	usr, ok := c.Get("userId").(*model.UserResponse)
	if ok {
		request.FriendId = request.UserId
		request.UserId = usr.Id.String()
	}

	if !ok {
		return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
			Code:    model.ErrResBadRequest.Code,
			Message: model.ErrResBadRequest.Message,
			Error:   err,
		})
	}

	if request.FriendId == request.UserId {
		return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
			Code:    model.ErrResBadRequest.Code,
			Message: "Cannot Delete Yourself",
			Error:   err,
		})
	}

	_, err = h.UseCase.RemoveFriend(c.Request().Context(), request.UserId, request.FriendId)

	if err != nil {

		if condition := errors.Is(err, model.ErrInvalidUserId); condition {

			return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
				Code:    model.ErrResBadRequest.Code,
				Message: model.ErrInvalidUserId.Error(),
				Error:   err,
			})
		}

		if condition := errors.Is(err, model.ErrUserNotFound); condition {
			return c.JSON(model.ErrResNotFound.Code, model.ResponseError{
				Code:    model.ErrResNotFound.Code,
				Message: model.ErrUserNotFound.Error(),
				Error:   err,
			})

		}

		if condition := errors.Is(err, model.ErrResNotFound.Error); condition {
			return c.JSON(model.ErrResNotFound.Code, model.ResponseError{
				Code:    model.ErrResNotFound.Code,
				Message: model.ErrUserNotFound.Error(),
				Error:   err,
			})

		}

		if condition := errors.Is(err, model.ErrNotFriend); condition {
			return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
				Code:    model.ErrResBadRequest.Code,
				Message: model.ErrNotFriend.Error(),
				Error:   err,
			})
		}

		return c.JSON(echo.ErrInternalServerError.Code, model.ResponseError{
			Code:    echo.ErrInternalServerError.Code,
			Message: echo.ErrInternalServerError.Error(),
			Error:   err,
		})
	}

	return c.JSON(http.StatusOK, model.Response[any]{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    make(map[string]interface{}),
	})
}
