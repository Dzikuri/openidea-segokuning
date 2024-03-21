package usecase

import (
	"context"

	"github.com/Dzikuri/openidea-segokuning/internal/model"
	"github.com/Dzikuri/openidea-segokuning/internal/repository"
	"github.com/rs/zerolog"
)

type UseCase interface {
	UserRegister(ctx context.Context, request *model.UserAuthRequest) (*model.UserAuthResponse, error)
	UserLogin(ctx context.Context, request *model.UserLoginRequest) (*model.UserAuthResponse, error)
	UserLinkEmail(ctx context.Context, request *model.UserLinkEmailRequest) (*model.UserResponse, error)
	UserLinkPhone(ctx context.Context, request *model.UserLinkPhoneRequest) (*model.UserResponse, error)
	GetUserByID(ctx context.Context, id string) (*model.UserResponse, int, error)
}

type useCase struct {
	Logger         zerolog.Logger
	UserRepository repository.RepositoryUser
}

func NewUseCase(logger zerolog.Logger, userRepository repository.RepositoryUser) UseCase {
	return &useCase{
		Logger:         logger,
		UserRepository: userRepository,
	}
}
