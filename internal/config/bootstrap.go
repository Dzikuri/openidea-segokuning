package config

import (
	"database/sql"

	"github.com/Dzikuri/openidea-segokuning/internal/delivery/handler"
	"github.com/Dzikuri/openidea-segokuning/internal/delivery/middleware"
	"github.com/Dzikuri/openidea-segokuning/internal/delivery/routes"
	"github.com/Dzikuri/openidea-segokuning/internal/repository"
	"github.com/Dzikuri/openidea-segokuning/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type BootstrapConfig struct {
	DB     *sql.DB
	App    *echo.Echo
	Logger *zerolog.Logger
}

func Bootstrap(config *BootstrapConfig) {

	userRepository := repository.NewUserRepository(config.DB)

	UseCase := usecase.NewUseCase(*config.Logger, userRepository)

	middleware := middleware.NewMiddleware(config.Logger, UseCase)

	handler := handler.NewHandler(UseCase, *config.Logger, middleware)

	routeConfig := routes.RoutesConfig{
		Echo:       config.App,
		Middleware: middleware,
		Handler:    *handler,
	}

	routeConfig.Setup()
}
