package usecase

import (
	"context"

	"github.com/Dzikuri/openidea-segokuning/internal/model"
)

type PostInterface interface {
	PostCreate(ctx context.Context, request *model.CreatePostRequest) (*model.PostResponse, error)
	PostCreateComment(ctx context.Context, request *model.CreatePostCommentRequest) (*model.PostCommentResponse, error)
	PostList(ctx context.Context, request *model.PostListRequest) (model.PaginateResponse[model.PostListResponse], error)
}

func (u *useCase) PostCreate(ctx context.Context, request *model.CreatePostRequest) (*model.PostResponse, error) {

	res, err := u.PostRepository.CreatePost(ctx, request)

	if err != nil {
		return nil, err
	}

	return res, nil

}

func (u *useCase) PostCreateComment(ctx context.Context, request *model.CreatePostCommentRequest) (*model.PostCommentResponse, error) {

	// NOTE Check Post is Exists
	_, err := u.PostRepository.FindPostById(ctx, request.PostId)
	if err != nil {
		return nil, err
	}

	res, err := u.PostRepository.CreatePostComment(ctx, request)

	if err != nil {
		return nil, err
	}

	return res, nil

}

func (u *useCase) PostList(ctx context.Context, request *model.PostListRequest) (model.PaginateResponse[model.PostListResponse], error) {

	res, meta, err := u.PostRepository.PostList(ctx, request)

	if err != nil {
		return model.PaginateResponse[model.PostListResponse]{
			Data: []model.PostListResponse{},
			Meta: model.MetaDataResponse{
				Total:  0,
				Limit:  request.Limit,
				Offset: request.Offset,
			},
			Message: "Ok",
		}, err
	}

	return model.PaginateResponse[model.PostListResponse]{
		Data:    res,
		Meta:    meta,
		Message: "Ok",
	}, err
}
