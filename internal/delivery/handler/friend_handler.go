package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

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

	params := c.QueryParams()

	var (
		limit, offset int
		err           error
	)

	queryParams := c.Request().URL.RawQuery
	if params.Get("limit") == "" {
		if strings.Contains(queryParams, "limit") {
			return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
				Code:    model.ErrResBadRequest.Code,
				Message: "Invalid query limit",
				Error:   errors.New("Invalid Limit value"),
			})
		} else {
			request.Limit = 5
		}
	} else {

		limit, err = strconv.Atoi(params.Get("limit"))
		if err != nil {
			return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
				Code:    model.ErrResBadRequest.Code,
				Message: "Invalid parsing query limit args to int",
				Error:   errors.New("Invalid limit value"),
			})
		}
		if limit < 0 {
			return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
				Code:    model.ErrResBadRequest.Code,
				Message: "Invalid query limit",
				Error:   errors.New("Invalid limit value"),
			})
		}
	}

	if params.Get("offset") == "" {
		if strings.Contains(queryParams, "offset") {
			return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
				Code:    model.ErrResBadRequest.Code,
				Message: "Invalid query offset",
				Error:   errors.New("Invalid offset value"),
			})
		} else {
			request.Offset = 0
		}
	} else {

		offset, err = strconv.Atoi(params.Get("offset"))
		if err != nil {
			return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
				Code:    model.ErrResBadRequest.Code,
				Message: "Invalid parsing query offset args to int",
				Error:   errors.New("Invalid offset value"),
			})
		}
		if offset < 0 {
			return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
				Code:    model.ErrResBadRequest.Code,
				Message: "Invalid query offset",
				Error:   errors.New("Invalid offset value"),
			})
		}
	}

	if params.Get("sortBy") == "" {
		if strings.Contains(queryParams, "sortBy") {
			return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
				Code:    model.ErrResBadRequest.Code,
				Message: "Invalid query sortBy",
				Error:   errors.New("Invalid sortBy value"),
			})
		} else {
			request.SortBy = "created_at"
		}
	} else if params.Get("sortBy") != "createdAt" && params.Get("sortBy") != "friendCount" {
		return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
			Code:    model.ErrResBadRequest.Code,
			Message: "Invalid query sortBy",
			Error:   errors.New("Invalid sortBy value"),
		})
	}

	if params.Get("orderBy") == "" {
		if strings.Contains(queryParams, "orderBy") {
			return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
				Code:    model.ErrResBadRequest.Code,
				Message: "Invalid query orderBy",
				Error:   errors.New("Invalid orderBy value"),
			})
		} else {
			request.OrderBy = "desc"
		}
	} else if condition := params.Get("orderBy") != "asc" && params.Get("orderBy") != "desc"; condition {
		return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
			Code:    model.ErrResBadRequest.Code,
			Message: "Invalid query orderBy",
			Error:   errors.New("Invalid orderBy value"),
		})
	}

	if params.Get("onlyFriend") == "" {
		if strings.Contains(queryParams, "onlyFriend") {
			return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
				Code:    model.ErrResBadRequest.Code,
				Message: "Invalid query onlyFriend",
				Error:   errors.New("Invalid onlyFriend value"),
			})
		} else {
			request.OnlyFriend = false
		}
	}

	err = c.Bind(&request)
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
