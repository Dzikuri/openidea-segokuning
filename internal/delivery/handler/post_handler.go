package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Dzikuri/openidea-segokuning/internal/model"
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreatePost(c echo.Context) error {

	var request model.CreatePostRequest
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
			Message: err.Error(),
			Error:   err,
		})
	}

	usr, ok := c.Get("userId").(*model.UserResponse)
	if ok {
		request.UserId = usr.Id.String()
	}

	result, err := h.UseCase.PostCreate(c.Request().Context(), &request)
	if err != nil {
		return c.JSON(echo.ErrInternalServerError.Code, model.ResponseError{
			Code:    echo.ErrInternalServerError.Code,
			Message: err.Error(),
			Error:   err,
		})
	}

	return c.JSON(http.StatusOK, model.Response[any]{
		Message: "Post created",
		Code:    http.StatusOK,
		Data:    result,
	})
}

func (h *Handler) GetPosts(c echo.Context) error {

	// Parse query parameters
	limit := c.Request().URL.Query().Get("limit")
	offset := c.Request().URL.Query().Get("offset")

	var request model.PostListRequest

	err := c.Bind(&request)
	if err != nil {
		return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
			Code:    model.ErrResBadRequest.Code,
			Message: model.ErrResBadRequest.Message,
			Error:   err,
		})
	}

	// Convert string parameters to integers
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		// Handle error (e.g., invalid parameter format)
		return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
			Code:    model.ErrResBadRequest.Code,
			Message: "Invalid Limit value",
			Error:   err,
		})
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		// Handle error (e.g., invalid parameter format)
		return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
			Code:    model.ErrResBadRequest.Code,
			Message: "Invalid offset value",
			Error:   err,
		})
	}

	request.Limit = limitInt
	request.Offset = offsetInt

	err = request.Validate()
	if err != nil {
		return c.JSON(model.ErrResBadRequest.Code, model.ResponseError{
			Code:    model.ErrResBadRequest.Code,
			Message: model.ErrResBadRequest.Message,
			Error:   err,
		})
	}

	res, err := h.UseCase.PostList(c.Request().Context(), &request)

	if err != nil {
		return c.JSON(echo.ErrInternalServerError.Code, model.ResponseError{
			Code:    echo.ErrInternalServerError.Code,
			Message: err.Error(),
			Error:   err,
		})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) CreatePostComment(c echo.Context) error {

	var request model.CreatePostCommentRequest
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
			Message: err.Error(),
			Error:   err,
		})
	}

	usr, ok := c.Get("userId").(*model.UserResponse)
	if ok {
		request.UserId = usr.Id.String()
	}

	result, err := h.UseCase.PostCreateComment(c.Request().Context(), &request)
	if err != nil {

		if errors.Is(err, model.ErrResNotFound.Error) {
			return c.JSON(echo.ErrNotFound.Code, model.ResponseError{
				Code:    echo.ErrNotFound.Code,
				Message: "Post Not Found",
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
		Message: "Post created",
		Code:    http.StatusOK,
		Data:    result,
	})

}
