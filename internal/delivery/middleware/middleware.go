package middleware

import (
	"github.com/Dzikuri/openidea-segokuning/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type Middleware interface {
	Auth() echo.MiddlewareFunc
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

func (m *middleware) Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}
