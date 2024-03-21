package usecase

import (
	"context"
	"fmt"

	"github.com/Dzikuri/openidea-segokuning/internal/helper"
	"github.com/Dzikuri/openidea-segokuning/internal/model"
	uuid "github.com/satori/go.uuid"
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

	var (
		exists bool
		result *model.UserResponse
		err    error
	)

	if requestAuth.CredentialType == model.Email {

		exists, result, err = u.UserRepository.FindByEmail(ctx, &requestAuth)
	}

	if requestAuth.CredentialType == model.Phone {
		exists, result, err = u.UserRepository.FindByPhone(ctx, &requestAuth)
	}

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

func (u *useCase) GetUserByID(ctx context.Context, id string) (*model.UserResponse, int, error) {
	user, code, err := u.UserRepository.FindById(ctx, id)
	if err != nil {
		return nil, code, err
	}
	return &model.UserResponse{
		Id:        user.Id,
		Email:     user.Email,
		Phone:     user.Phone,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, code, nil
}

func (u *useCase) UserLinkEmail(ctx context.Context, request *model.UserLinkEmailRequest) (*model.UserResponse, error) {

	fmt.Println("request: ")
	helper.LogPretty(request)
	exists, checkEmail, err := u.UserRepository.FindByEmail(ctx, &model.UserAuthRequest{
		CredentialType:  "email",
		CredentialValue: request.Email,
	})

	if checkEmail.Id != uuid.Nil && request.Id == checkEmail.Id {
		return nil, model.ErrLinkEmailExists
	}

	if checkEmail.Id != uuid.Nil {
		return nil, model.ErrUserAlreadyExists
	}

	fmt.Println("exists : ", exists)
	helper.LogPretty(checkEmail)
	fmt.Println("err : ", err)

	return nil, nil
}

func (u *useCase) UserLinkPhone(ctx context.Context, request *model.UserLinkPhoneRequest) (*model.UserResponse, error) {
	return nil, nil
}
