package usecase

import (
	"context"
	"errors"

	"github.com/Dzikuri/openidea-segokuning/internal/model"
)

type FriendInterface interface {
	// IsFriend(ctx context.Context, userID string, friendID string) (bool, error)
	AddFriend(ctx context.Context, userID string, friendID string) (*model.FriendResponse, error)
	RemoveFriend(ctx context.Context, userID string, friendID string) (*model.FriendResponse, error)
	GetFriendList(ctx context.Context, request model.GetFriendListRequest) (model.PaginateResponse[model.FriendResponse], error)
}

func (u *useCase) AddFriend(ctx context.Context, userID string, friendID string) (*model.FriendResponse, error) {
	// Check User Exists
	_, exists, err := u.FriendRepository.CheckAlreadyFriend(ctx, userID, friendID)
	if err != nil {
		if errors.Is(err, model.ErrResNotFound.Error) {
			return nil, model.ErrResNotFound.Error
		}

		return nil, err
	}

	if exists == 1 {
		return nil, model.ErrAlreadyBeFriend
	}

	result, err := u.FriendRepository.AddFriend(ctx, model.FriendRequest{
		UserId:   userID,
		FriendId: friendID,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *useCase) RemoveFriend(ctx context.Context, userID string, friendID string) (*model.FriendResponse, error) {
	// Check User Exists
	_, exists, err := u.FriendRepository.CheckAlreadyFriend(ctx, userID, friendID)
	if err != nil {
		if errors.Is(err, model.ErrResNotFound.Error) {
			return nil, model.ErrResNotFound.Error
		}

		return nil, err
	}

	if exists == 0 {
		return nil, model.ErrNotFriend
	}

	result, err := u.FriendRepository.RemoveFriend(ctx, userID, friendID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *useCase) GetFriendList(ctx context.Context, request model.GetFriendListRequest) (model.PaginateResponse[model.FriendResponse], error) {

	if request.Limit == 0 {
		request.Limit = 5
	}

	if request.SortBy == "" {
		request.SortBy = "createdAt"
	}

	if request.OrderBy == "" {
		request.OrderBy = "desc"
	}

	result, meta, err := u.FriendRepository.FindAllFriend(ctx, request)
	if err != nil {
		return model.PaginateResponse[model.FriendResponse]{
			Data: []model.FriendResponse{},
			Meta: model.MetaDataResponse{
				Total:  0,
				Limit:  request.Limit,
				Offset: request.Offset,
			},
			Message: "Ok",
		}, err
	}

	return model.PaginateResponse[model.FriendResponse]{
		Data:    result,
		Meta:    meta,
		Message: "Ok",
	}, nil
}
