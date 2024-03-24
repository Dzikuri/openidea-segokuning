package handler

import (
	"github.com/Dzikuri/openidea-segokuning/internal/delivery/middleware"
	"github.com/Dzikuri/openidea-segokuning/internal/usecase"
	"github.com/rs/zerolog"
)

type Handler struct {
	UseCase    usecase.UseCase
	Logger     zerolog.Logger
	Middleware middleware.Middleware
}

func NewHandler(usecase usecase.UseCase, logger zerolog.Logger, middleware middleware.Middleware) *Handler {
	return &Handler{
		UseCase:    usecase,
		Logger:     logger,
		Middleware: middleware,
	}
}
