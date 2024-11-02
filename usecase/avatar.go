package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/todennus/shared/enumdef"
	"github.com/todennus/shared/errordef"
	"github.com/todennus/shared/scopedef"
	"github.com/todennus/shared/xcontext"
	"github.com/todennus/user-service/usecase/abstraction"
	"github.com/todennus/user-service/usecase/dto"
	"github.com/todennus/x/xerror"
)

type AvatarUsecase struct {
	avatarDomain abstraction.AvatarDomain

	avatarPolicyRepo abstraction.AvatarPolicySessionRepository
	fileRepo         abstraction.FileRepository
	userRepo         abstraction.UserRepository
}

func NewAvatarUsecase(
	avatarDomain abstraction.AvatarDomain,
	avatarPolicyRepo abstraction.AvatarPolicySessionRepository,
	fileRepo abstraction.FileRepository,
	userRepo abstraction.UserRepository,
) *AvatarUsecase {
	return &AvatarUsecase{
		avatarDomain:     avatarDomain,
		avatarPolicyRepo: avatarPolicyRepo,
		fileRepo:         fileRepo,
		userRepo:         userRepo,
	}
}

func (usecase *AvatarUsecase) GetPolicyToken(
	ctx context.Context,
	req *dto.AvatarGetPolicyTokenRequest,
) (*dto.AvatarGetPolicyTokenResponse, error) {
	if scopedef.Eval(xcontext.Scope(ctx)).RequireAnyUser(scopedef.UserUpdateUserAvatar).IsUnsatisfied() {
		return nil, xerror.Enrich(errordef.ErrForbidden, "insufficient scope")
	}

	if xcontext.RequestSubjectType(ctx) != enumdef.SubjectUser {
		return nil, xerror.Enrich(errordef.ErrForbidden, "require a user subject")
	}

	policy := usecase.avatarDomain.GetPolicy(xcontext.RequestSubjectID(ctx))
	if err := usecase.avatarPolicyRepo.Store(ctx, policy); err != nil {
		return nil, errordef.ErrServer.Hide(err, "failed-to-store-avatar-policy")
	}

	return dto.NewAvatarGetPolicyTokenResponse(policy.PolicyToken), nil
}

func (usecase *AvatarUsecase) ValidatePolicyToken(
	ctx context.Context,
	req *dto.AvatarValidatePolicyTokenRequest,
) (*dto.AvatarValidatePolicyTokenResponse, error) {
	if scopedef.Eval(xcontext.Scope(ctx)).RequireAdmin(scopedef.AdminValidateFilePolicy).IsUnsatisfied() {
		return nil, xerror.Enrich(errordef.ErrForbidden, "insufficient scope")
	}

	policy, err := usecase.avatarPolicyRepo.Load(ctx, req.PolicyToken)
	if err != nil {
		return nil, xerror.Enrich(errordef.ErrForbidden, "invalid policy token")
	}

	if policy.ExpiresAt.Before(time.Now()) {
		return nil, xerror.Enrich(errordef.ErrForbidden, "expired policy token")
	}

	if err := usecase.avatarPolicyRepo.Delete(ctx, req.PolicyToken); err != nil {
		xcontext.Logger(ctx).Warn("failed-to-delete-avatar-policy-token", "err", err)
	}

	return dto.NewAvatarValidatePolicyTokenResponse(policy), nil
}

func (usecase *AvatarUsecase) Update(
	ctx context.Context,
	req *dto.AvatarUpdateRequest,
) (*dto.AvatarUpdateResponse, error) {
	if scopedef.Eval(xcontext.Scope(ctx)).RequireAnyUser(scopedef.UserUpdateUserAvatar).IsUnsatisfied() {
		return nil, xerror.Enrich(errordef.ErrForbidden, "insufficient scope")
	}

	if xcontext.RequestSubjectType(ctx) != enumdef.SubjectUser {
		return nil, xerror.Enrich(errordef.ErrForbidden, "require a user subject")
	}

	userID, err := usecase.fileRepo.ValidateTemporaryFile(ctx, req.TemporaryFileToken)
	if err != nil {
		if errors.Is(err, errordef.ErrRequestInvalid) {
			return nil, xerror.Enrich(errordef.ErrForbidden, "invalid session token")
		}

		return nil, errordef.ErrServer.Hide(err, "failed-to-validate-image-session")
	}

	if userID != xcontext.RequestSubjectID(ctx) {
		if err := usecase.fileRepo.DeleteTemporary(ctx, req.TemporaryFileToken); err != nil {
			xcontext.Logger(ctx).Warn("failed-to-reject-image-session", "err", err)
		}

		return nil, xerror.Enrich(errordef.ErrForbidden, "the user doesn't have the permission to use this token")
	}

	avatarURL, err := usecase.fileRepo.SaveToPersistent(ctx, req.TemporaryFileToken)
	if err != nil {
		return nil, errordef.ErrServer.Hide(err, "failed-to-accept-image-session")
	}

	if err := usecase.userRepo.UpdateAvatarByID(ctx, userID, avatarURL); err != nil {
		return nil, errordef.ErrServer.Hide(err, "failed-to-update-avatar", "avatar_url", avatarURL)
	}

	return dto.NewAvatarUpdateResponse(avatarURL), nil
}
