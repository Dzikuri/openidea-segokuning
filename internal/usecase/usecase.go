package usecase

import (
	"github.com/Dzikuri/openidea-segokuning/internal/repository"
	"github.com/rs/zerolog"
)

type UseCase interface {
	UserInterface
	FriendInterface
}

type useCase struct {
	Logger           zerolog.Logger
	UserRepository   repository.RepositoryUser
	FriendRepository repository.RepositoryFriend
}

func NewUseCase(logger zerolog.Logger, userRepository repository.RepositoryUser, friendRepository repository.RepositoryFriend) UseCase {
	return &useCase{
		Logger:           logger,
		UserRepository:   userRepository,
		FriendRepository: friendRepository,
	}
}
