package config

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func NewEcho(logger *zerolog.Logger) *echo.Echo {
	// Init echo server
	e := echo.New()

	e.HideBanner = true

	// e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("URI", v.URI).
				Str("method", c.Request().Method).
				Int("status", v.Status).
				Msg("request")
			return nil
		},
	}))

	return e
}
