package usecase

import (
	"context"
	"slices"

	"github.com/todennus/shared/errordef"
	"github.com/todennus/shared/scopedef"
	"github.com/todennus/shared/tokendef"
	"github.com/todennus/shared/xcontext"
	"github.com/todennus/user-service/usecase/abstraction"
	"github.com/todennus/user-service/usecase/dto"
	"github.com/todennus/x/token"
	"github.com/todennus/x/xerror"
	"github.com/xybor-x/snowflake"
)

type AvatarUsecase struct {
	tokenEngine token.Engine

	avatarDomain abstraction.AvatarDomain

	fileRepo abstraction.FileRepository
	userRepo abstraction.UserRepository
}

func NewAvatarUsecase(
	tokenEngine token.Engine,
	avatarDomain abstraction.AvatarDomain,
	fileRepo abstraction.FileRepository,
	userRepo abstraction.UserRepository,
) *AvatarUsecase {
	return &AvatarUsecase{
		tokenEngine:  tokenEngine,
		avatarDomain: avatarDomain,
		fileRepo:     fileRepo,
		userRepo:     userRepo,
	}
}

func (usecase *AvatarUsecase) GetUploadToken(
	ctx context.Context,
	req *dto.AvatarGetUploadTokenRequest,
) (*dto.AvatarGetUploadTokenResponse, error) {
	if scopedef.Eval(xcontext.Scope(ctx)).RequireAnyUser(scopedef.UserUpdateUserAvatar).IsUnsatisfied() {
		return nil, xerror.Enrich(errordef.ErrForbidden, "insufficient scope")
	}

	if req.UserID != xcontext.RequestSubjectID(ctx) {
		return nil, xerror.Enrich(errordef.ErrForbidden, "permission denied")
	}

	uploadToken, err := usecase.fileRepo.RegisterUpload(
		ctx, usecase.avatarDomain.GetPolicy(xcontext.RequestSubjectID(ctx)))
	if err != nil {
		return nil, errordef.ErrServer.Hide(err, "failed-to-register-upload-token")
	}

	return dto.NewAvatarGetUploadTokenResponse(uploadToken), nil
}

func (usecase *AvatarUsecase) Update(
	ctx context.Context,
	req *dto.AvatarUpdateRequest,
) (*dto.AvatarUpdateResponse, error) {
	if scopedef.Eval(xcontext.Scope(ctx)).RequireAnyUser(scopedef.UserUpdateUserAvatar).IsUnsatisfied() {
		return nil, xerror.Enrich(errordef.ErrForbidden, "insufficient scope")
	}

	if req.UserID != xcontext.RequestSubjectID(ctx) {
		return nil, xerror.Enrich(errordef.ErrForbidden, "permission denied")
	}

	fileToken := &tokendef.FileToken{}
	if err := usecase.tokenEngine.Validate(ctx, req.FileToken, fileToken); err != nil {
		return nil, xerror.Enrich(errordef.ErrRequestInvalid, "invalid token").Hide(err, "failed-to-parse-token")
	}

	policy := usecase.avatarDomain.GetPolicy(xcontext.RequestSubjectID(ctx))

	userID := fileToken.SnowflakeUserID()
	if userID != policy.UserID {
		return nil, xerror.Enrich(errordef.ErrForbidden, "the user doesn't have the permission to use this token")
	}

	if !slices.Contains(policy.AllowedTypes, fileToken.Type) {
		return nil, xerror.Enrich(errordef.ErrFileMismatchedType,
			"require %s, but got %s", policy.AllowedTypes, fileToken.Type)
	}

	if fileToken.Size > int(policy.MaxSize) {
		return nil, xerror.Enrich(errordef.ErrRequestInvalid,
			"the image size is limited at %d bytes, but got %d", policy.MaxSize, fileToken.Size)
	}

	ctx = xcontext.WithDBTransaction(ctx)
	defer xcontext.DBCommit(ctx)

	currentAvatar, err := usecase.userRepo.GetAvatarByID(ctx, userID)
	if err != nil {
		return nil, errordef.ErrServer.Hide(err, "failed-to-get-current-avatar")
	}

	newAvatar := fileToken.SnowflakeOwnershipID()
	if err := usecase.userRepo.UpdateAvatarByID(ctx, userID, newAvatar); err != nil {
		ctx = xcontext.DBRollback(ctx)
		return nil, errordef.ErrServer.Hide(err, "failed-to-update-avatar",
			"user_id", userID, "ownership_id", newAvatar)
	}

	incOwnershipID := []snowflake.ID{newAvatar}
	decOwnershipID := []snowflake.ID{}
	if currentAvatar != 0 {
		decOwnershipID = append(decOwnershipID, currentAvatar)
	}

	if err := usecase.fileRepo.ChangeRefcount(ctx, incOwnershipID, decOwnershipID); err != nil {
		ctx = xcontext.DBRollback(ctx)
		return nil, errordef.ErrServer.Hide(err, "failed-to-change-avatar-ref-count",
			"user_id", userID, "inc", incOwnershipID, "dec", decOwnershipID)
	}

	return dto.NewAvatarUpdateResponse(), nil
}
