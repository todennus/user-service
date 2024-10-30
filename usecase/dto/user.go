package dto

import (
	"context"

	"github.com/todennus/user-service/domain"
	"github.com/todennus/user-service/usecase/dto/resource"
	"github.com/xybor-x/snowflake"
)

type UserRegisterRequest struct {
	Username string
	Password string
}

type UserRegisterResponse struct {
	User *resource.User
}

func NewUserRegisterResponse(user *domain.User) *UserRegisterResponse {
	return &UserRegisterResponse{
		User: resource.NewUserWithoutFilter(user),
	}
}

type UserRegisterFirstRequest struct {
	Username string
	Password string
}

type UserRegisterFirstResponse struct {
	User *resource.User
}

func NewUserRegisterFirstResponse(user *domain.User) *UserRegisterFirstResponse {
	return &UserRegisterFirstResponse{
		User: resource.NewUserWithoutFilter(user),
	}
}

type UserGetByIDRequest struct {
	UserID snowflake.ID
}

type UserGetByIDResponse struct {
	User *resource.User
}

func NewUserGetByIDResponse(ctx context.Context, user *domain.User) *UserGetByIDResponse {
	return &UserGetByIDResponse{
		User: resource.NewUserWithFilter(ctx, user),
	}
}

type UserGetByUsernameRequest struct {
	Username string
}

type UserGetByUsernameResponse struct {
	User *resource.User
}

func NewUserGetByUsernameResponse(ctx context.Context, user *domain.User) *UserGetByUsernameResponse {
	return &UserGetByUsernameResponse{
		User: resource.NewUserWithFilter(ctx, user),
	}
}

type UserValidateCredentialsRequest struct {
	Username string
	Password string
}

type UserValidateCredentialsResponse struct {
	User *resource.User
}

func NewUserValidateCredentialsResponse(user *domain.User) *UserValidateCredentialsResponse {
	return &UserValidateCredentialsResponse{
		User: resource.NewUserWithoutFilter(user),
	}
}
