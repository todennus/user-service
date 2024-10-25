package usecase

import (
	"context"
	"errors"

	"github.com/todennus/shared/enumdef"
	"github.com/todennus/shared/errordef"
	"github.com/todennus/user-service/domain"
	"github.com/todennus/user-service/usecase/abstraction"
	"github.com/todennus/user-service/usecase/dto"
	"github.com/todennus/x/lock"
	"github.com/todennus/x/xcontext"
	"github.com/todennus/x/xerror"
)

type UserUsecase struct {
	adminLocker       lock.Locker
	shouldCreateAdmin bool

	userDomain abstraction.UserDomain
	userRepo   abstraction.UserRepository
}

func NewUserUsecase(
	locker lock.Locker,
	userRepo abstraction.UserRepository,
	userDomain abstraction.UserDomain,
) *UserUsecase {
	return &UserUsecase{
		adminLocker:       locker,
		shouldCreateAdmin: true,
		userRepo:          userRepo,
		userDomain:        userDomain,
	}
}

func (uc *UserUsecase) Register(
	ctx context.Context,
	req *dto.UserRegisterRequest,
) (*dto.UserRegisterResponse, error) {
	_, err := uc.userRepo.GetByUsername(ctx, req.Username)
	if err == nil {
		return nil, xerror.Enrich(errordef.ErrDuplicated, "username %s has already existed", req.Username)
	}

	if !errors.Is(err, errordef.ErrNotFound) {
		return nil, errordef.ErrServer.Hide(err, "failed-to-get-user")
	}

	user, err := uc.userDomain.Create(req.Username, req.Password)
	if err != nil {
		return nil, errordef.Domain.Event(err, "failed-to-new-user").Enrich(errordef.ErrRequestInvalid).Error()
	}

	shouldCreateUser, err := uc.createAdmin(ctx, user)
	if err != nil {
		return nil, err
	}

	if shouldCreateUser {
		if err = uc.userRepo.Create(ctx, user); err != nil {
			return nil, errordef.ErrServer.Hide(err, "failed-to-create-user")
		}
	}

	return dto.NewUserRegisterResponse(ctx, user), nil
}

func (usecase *UserUsecase) GetByID(
	ctx context.Context,
	req *dto.UserGetByIDRequest,
) (*dto.UserGetByIDResponse, error) {
	if req.UserID == 0 {
		return nil, xerror.Enrich(errordef.ErrRequestInvalid, "require user id")
	}

	user, err := usecase.userRepo.GetByID(ctx, req.UserID.Int64())
	if err != nil {
		if errors.Is(err, errordef.ErrNotFound) {
			return nil, xerror.Enrich(errordef.ErrNotFound, "not found user with id %d", req.UserID)
		}

		return nil, errordef.ErrServer.Hide(err, "failed-to-get-user", "uid", req.UserID)
	}

	return dto.NewUserGetByIDResponse(ctx, user), nil
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

	return dto.NewUserGetByUsernameResponse(ctx, user), nil
}

func (usecase *UserUsecase) ValidateCredentials(
	ctx context.Context,
	req *dto.UserValidateCredentialsRequest,
) (*dto.UserValidateCredentialsResponse, error) {
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
		return nil, errordef.Domain.Event(err, "failed-to-validate-user-credentials").
			EnrichWith(errordef.ErrCredentialsInvalid, "invalid username or password").
			Error()
	}

	ctx = xcontext.WithRequestUserID(ctx, user.ID)
	return dto.NewUserValidateCredentialsResponse(user), nil
}

func (uc *UserUsecase) createAdmin(
	ctx context.Context,
	user *domain.User,
) (bool, error) {
	if !uc.shouldCreateAdmin {
		return true, nil
	}

	shouldCreateUser := true
	xcontext.Logger(ctx).Info("check-create-admin")

	err := lock.Func(uc.adminLocker, ctx, func() error {
		adminCount, err := uc.userRepo.CountByRole(ctx, enumdef.UserRoleAdmin)
		if err != nil {
			return errordef.ErrServer.Hide(err, "failed-to-count-by-admin-role")
		}

		if adminCount > 0 {
			xcontext.Logger(ctx).Info("cancel-create-admin")
			uc.shouldCreateAdmin = false
			return nil
		}

		xcontext.Logger(ctx).Info("create-admin", "username", user.Username, "uid", user.ID)

		user.Role = enumdef.UserRoleAdmin
		err = uc.userRepo.Create(ctx, user)
		if err != nil {
			return errordef.ErrServer.Hide(err, "failed-to-create-admin-user")
		}

		shouldCreateUser = false
		uc.shouldCreateAdmin = false
		return nil
	})

	return shouldCreateUser, err
}
