package usecase

import (
	"context"

	"github.com/Dzikuri/openidea-segokuning/internal/helper"
	"github.com/Dzikuri/openidea-segokuning/internal/model"
)

func (u *useCase) UserRegister(ctx context.Context, request *model.UserAuthRequest) (*model.UserAuthResponse, error) {
	// hash password before insert to database

	hashedPassword, err := helper.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	// insert to database

	request.Password = hashedPassword

	result, err := u.UserRepository.Register(ctx, request)
	if err != nil {
		return nil, err
	}

	// Generate JWT
	token, err := helper.JwtGenerateToken(result)
	if err != nil {
		return nil, err
	}

	// return response

	return &model.UserAuthResponse{
		Phone:       result.Phone,
		Email:       result.Email,
		Name:        result.Name,
		AccessToken: token,
	}, nil
}

func (u *useCase) UserLogin(ctx context.Context, request *model.UserLoginRequest) (*model.UserAuthResponse, error) {

	var requestAuth model.UserAuthRequest

	requestAuth.CredentialType = request.CredentialType
	requestAuth.CredentialValue = request.CredentialValue
	requestAuth.Password = request.Password

	exists, result, err := u.UserRepository.FindByEmail(ctx, &requestAuth)

	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, model.ErrUserNotFound
	}

	err = helper.ComparePassword(result.Password, request.Password)

	if err != nil {
		return nil, model.ErrPasswordNotMatch
	}

	// Generate JWT
	token, err := helper.JwtGenerateToken(result)
	if err != nil {
		return nil, err
	}
	// return response
	return &model.UserAuthResponse{
		Phone:       result.Phone,
		Email:       result.Email,
		Name:        result.Name,
		AccessToken: token,
	}, nil
}
