package wiring

import (
	"context"
	"time"

	"github.com/todennus/shared/config"
	"github.com/todennus/user-service/adapter/abstraction"
	"github.com/todennus/user-service/usecase"
	"github.com/todennus/x/lock"
)

type Usecases struct {
	abstraction.UserUsecase
	abstraction.AvatarUsecase
}

func InitializeUsecases(
	ctx context.Context,
	config *config.Config,
	infras *Infras,
	domains *Domains,
	repositories *Repositories,
) (*Usecases, error) {
	uc := &Usecases{}

	uc.UserUsecase = usecase.NewUserUsecase(
		lock.NewRedisLock(infras.Redis, "user-lock", 10*time.Second),
		time.Duration(config.Variable.User.AvatarPresignedURLExpiration)*time.Second,
		domains.UserDomain,
		repositories.UserRepository,
		repositories.FileRepository,
	)

	uc.AvatarUsecase = usecase.NewAvatarUsecase(
		config.TokenEngine,
		domains.AvatarDomain,
		repositories.FileRepository,
		repositories.UserRepository,
	)

	return uc, nil
}
