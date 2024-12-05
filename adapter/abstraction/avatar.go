package abstraction

import (
	"context"

	"github.com/todennus/user-service/usecase/dto"
)

type AvatarUsecase interface {
	GetUploadToken(context.Context, *dto.AvatarGetUploadTokenRequest) (*dto.AvatarGetUploadTokenResponse, error)
	Update(context.Context, *dto.AvatarUpdateRequest) (*dto.AvatarUpdateResponse, error)
}
