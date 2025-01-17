package abstraction

import (
	"context"

	"github.com/todennus/user-service/usecase/dto"
)

type UserUsecase interface {
	Register(ctx context.Context, req *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error)
	RegisterFirst(ctx context.Context, req *dto.UserRegisterFirstRequest) (*dto.UserRegisterFirstResponse, error)
	GetByID(ctx context.Context, req *dto.UserGetByIDRequest) (*dto.UserGetByIDResponse, error)
	GetByUsername(ctx context.Context, req *dto.UserGetByUsernameRequest) (*dto.UserGetByUsernameResponse, error)
	ValidateCredentials(
		ctx context.Context,
		req *dto.UserValidateCredentialsRequest,
	) (*dto.UserValidateCredentialsResponse, error)
}
