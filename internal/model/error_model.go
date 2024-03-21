package model

import (
	"errors"

	"github.com/labstack/echo/v4"
)

var (
	ErrResForbidden     = ResponseError{Code: echo.ErrForbidden.Code, Message: "Forbidden", Error: errors.New("Forbidden")}
	ErrResBadRequest    = ResponseError{Code: echo.ErrBadRequest.Code, Message: "Bad Request", Error: errors.New("Bad Request")}
	ErrResUnauthorized  = ResponseError{Code: echo.ErrUnauthorized.Code, Message: "Unauthorized", Error: errors.New("Unauthorized")}
	ErrResRequiredField = ResponseError{Code: echo.ErrBadRequest.Code, Message: "Required Field", Error: errors.New("Required Field")}
	ErrResNotFound      = ResponseError{Code: echo.ErrNotFound.Code, Message: "Not Found", Error: errors.New("Not Found")}
	ErrResNoRecord      = ResponseError{Code: echo.ErrNotFound.Code, Message: "No Record", Error: errors.New("No Record")}
)

var (
	ErrUserAlreadyExists = errors.New("User already exists")
	ErrUserNotFound      = errors.New("User not found")
	ErrUnauthorize       = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrPasswordNotMatch  = errors.New("password not match")
)
