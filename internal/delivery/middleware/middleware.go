package middleware

import (
	"net/http"
	"strings"

	"github.com/Dzikuri/openidea-segokuning/internal/helper"
	"github.com/Dzikuri/openidea-segokuning/internal/model"
	"github.com/Dzikuri/openidea-segokuning/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Middleware interface {
	Authentication(isThrowError bool) func(next echo.HandlerFunc) echo.HandlerFunc
}

type middleware struct {
	Logger  *zerolog.Logger
	UseCase usecase.UseCase
}

func NewMiddleware(logger *zerolog.Logger, useCase usecase.UseCase) Middleware {
	return &middleware{
		Logger:  logger,
		UseCase: useCase,
	}
}

func (m *middleware) Authentication(isThrowError bool) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			m.Logger.Info().Msg("Authentication")
			token := GetBearerTokenFromRequest(c.Request())

			if token == "" && isThrowError {
				return c.JSON(http.StatusUnauthorized, model.ResponseError{
					Code:    http.StatusUnauthorized,
					Message: model.ErrUnauthorize.Error(),
					// Error:   echo.ErrUnauthorized,
				})
			}

			if token != "" {
				claims := &helper.JwtCustomClaims{}
				err := helper.VerifyJwt(token, claims, helper.JwtSecret())
				if err != nil {
					return c.JSON(http.StatusUnauthorized, model.ResponseError{
						Code:    http.StatusUnauthorized,
						Message: errors.Cause(err).Error(),
						Error:   err,
					})
				}

				usr, code, err := m.UseCase.GetUserByID(c.Request().Context(), claims.Id)
				if err != nil {
					return c.JSON(code, model.ResponseError{
						Code:    code,
						Message: err.Error(),
					})
				}
				c.Set("userId", usr)
			}

			return next(c)
		}
	}
}

func GetBearerTokenFromRequest(r *http.Request) string {
	// From query.
	query := r.URL.Query().Get("jwt")
	if query != "" {
		return query
	}

	// From header.
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToLower(bearer[0:6]) == "bearer" {
		return bearer[7:]
	}

	// From cookie.
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return ""
	}

	return cookie.Value
}
