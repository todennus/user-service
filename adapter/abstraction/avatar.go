package abstraction

import (
	"context"

	"github.com/todennus/user-service/usecase/dto"
)

type AvatarUsecase interface {
	GetPolicyToken(context.Context, *dto.AvatarGetPolicyTokenRequest) (*dto.AvatarGetPolicyTokenResponse, error)
	ValidatePolicyToken(context.Context, *dto.AvatarValidatePolicyTokenRequest) (*dto.AvatarValidatePolicyTokenResponse, error)
	Update(context.Context, *dto.AvatarUpdateRequest) (*dto.AvatarUpdateResponse, error)
}
