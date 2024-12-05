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
	"github.com/todennus/x/lock"
	"github.com/todennus/x/xerror"
)

type UserUsecase struct {
	adminLocker       lock.Locker
	shouldCreateAdmin bool

	avatarPresignedURLExpiration time.Duration

	userDomain abstraction.UserDomain
	userRepo   abstraction.UserRepository
	fileRepo   abstraction.FileRepository
}

func NewUserUsecase(
	locker lock.Locker,
	avatarPresignedURLExpiration time.Duration,
	userDomain abstraction.UserDomain,
	userRepo abstraction.UserRepository,
	fileRepo abstraction.FileRepository,
) *UserUsecase {
	return &UserUsecase{
		adminLocker:                  locker,
		avatarPresignedURLExpiration: avatarPresignedURLExpiration,
		shouldCreateAdmin:            true,
		userRepo:                     userRepo,
		userDomain:                   userDomain,
		fileRepo:                     fileRepo,
	}
}

func (uc *UserUsecase) Register(
	ctx context.Context,
	req *dto.UserRegisterRequest,
) (*dto.UserRegisterResponse, error) {
	if scopedef.Eval(xcontext.Scope(ctx)).RequireAdmin(scopedef.AdminCreateUser).IsUnsatisfied() {
		return nil, xerror.Enrich(errordef.ErrForbidden, "insufficient scope")
	}

	user, err := uc.userDomain.New(req.Username, req.Password)
	if err != nil {
		return nil, errordef.DomainWrapper.Event(err, "failed-to-new-user").
			Enrich(errordef.ErrRequestInvalid).Error()
	}

	if err = uc.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, errordef.ErrDuplicated) {
			return nil, xerror.Enrich(errordef.ErrDuplicated, "username %s has already existed", req.Username)
		}

		return nil, errordef.ErrServer.Hide(err, "failed-to-create-user")
	}

	return dto.NewUserRegisterResponse(user), nil
}

func (uc *UserUsecase) RegisterFirst(
	ctx context.Context,
	req *dto.UserRegisterFirstRequest,
) (*dto.UserRegisterFirstResponse, error) {
	if !uc.shouldCreateAdmin {
		return nil, xerror.Enrich(errordef.ErrNotFound, "this api is only openned for creating the first user")
	}

	if err := uc.adminLocker.Lock(ctx); err != nil {
		return nil, errordef.ErrServer.Hide(err, "failed-to-lock-first-client-flow")
	}
	defer uc.adminLocker.Unlock(ctx)

	count, err := uc.userRepo.CountByRole(ctx, enumdef.UserRoleAdmin)
	if err != nil {
		return nil, errordef.ErrServer.Hide(err, "failed-to-count-client")
	}

	if count > 0 {
		uc.shouldCreateAdmin = false
		return nil, xerror.Enrich(errordef.ErrNotFound, "this api is only openned for creating the first user")
	}

	user, err := uc.userDomain.NewFirst(req.Username, req.Password)
	if err != nil {
		return nil, errordef.DomainWrapper.Event(err, "failed-to-new-user").
			Enrich(errordef.ErrRequestInvalid).Error()
	}

	if err = uc.userRepo.Create(ctx, user); err != nil {
		return nil, errordef.ErrServer.Hide(err, "failed-to-create-first-user")
	}

	uc.shouldCreateAdmin = false
	return dto.NewUserRegisterFirstResponse(user), nil
}

func (usecase *UserUsecase) GetByID(
	ctx context.Context,
	req *dto.UserGetByIDRequest,
) (*dto.UserGetByIDResponse, error) {
	if req.UserID == 0 {
		return nil, xerror.Enrich(errordef.ErrRequestInvalid, "require user id")
	}

	user, err := usecase.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, errordef.ErrNotFound) {
			return nil, xerror.Enrich(errordef.ErrNotFound, "not found user with id %d", req.UserID)
		}

		return nil, errordef.ErrServer.Hide(err, "failed-to-get-user", "uid", req.UserID)
	}

	var avatarURL string
	if user.Avatar != 0 {
		avatarURL, err = usecase.fileRepo.CreatePresignedURL(ctx, user.Avatar, usecase.avatarPresignedURLExpiration)
		if err != nil {
			return nil, errordef.ErrServer.Hide(err, "failed-to-get-presigned-url", "avatar", user.Avatar)
		}
	}

	return dto.NewUserGetByIDResponse(ctx, user, avatarURL), nil
}

func (usecase *UserUsecase) GetByUsername(
	ctx context.Context,
	req *dto.UserGetByUsernameRequest,
) (*dto.UserGetByUsernameResponse, error) {
	if req.Username == "" {
		return nil, xerror.Enrich(errordef.ErrRequestInvalid, "require username")
	}

	user, err := usecase.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, errordef.ErrNotFound) {
			return nil, xerror.Enrich(errordef.ErrNotFound, "not found user with username %s", req.Username)
		}

		return nil, errordef.ErrServer.Hide(err, "failed-to-get-user", "username", req.Username)
	}

	var avatarURL string
	if user.Avatar != 0 {
		avatarURL, err = usecase.fileRepo.CreatePresignedURL(ctx, user.Avatar, usecase.avatarPresignedURLExpiration)
		if err != nil {
			return nil, errordef.ErrServer.Hide(err, "failed-to-get-presigned-url", "avatar", user.Avatar)
		}
	}

	return dto.NewUserGetByUsernameResponse(ctx, user, avatarURL), nil
}

func (usecase *UserUsecase) ValidateCredentials(
	ctx context.Context,
	req *dto.UserValidateCredentialsRequest,
) (*dto.UserValidateCredentialsResponse, error) {
	if scopedef.Eval(xcontext.Scope(ctx)).RequireAdmin(scopedef.AdminValidateUser).IsUnsatisfied() {
		return nil, xerror.Enrich(errordef.ErrForbidden, "insufficient scope")
	}

	if req.Username == "" {
		return nil, xerror.Enrich(errordef.ErrRequestInvalid, "require username")
	}

	user, err := usecase.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, errordef.ErrNotFound) {
			return nil, xerror.Enrich(errordef.ErrCredentialsInvalid, "invalid username or password")
		}

		return nil, errordef.ErrServer.Hide(err, "failed-to-get-user", "username", req.Username)
	}

	if err := usecase.userDomain.Validate(user.HashedPass, req.Password); err != nil {
		return nil, errordef.DomainWrapper.Event(err, "failed-to-validate-user-credentials").
			EnrichWith(errordef.ErrCredentialsInvalid, "invalid username or password").
			Error()
	}

	ctx = xcontext.WithRequestSubjectID(ctx, user.ID)
	return dto.NewUserValidateCredentialsResponse(user), nil
}
