package wiring

import (
	"context"
	"time"

	"github.com/todennus/user-service/adapter/abstraction"
	"github.com/todennus/user-service/usecase"
	"github.com/todennus/x/lock"
)

type Usecases struct {
	abstraction.UserUsecase
}

func InitializeUsecases(
	ctx context.Context,
	infras *Infras,
	domains *Domains,
	repositories *Repositories,
) (*Usecases, error) {
	uc := &Usecases{}

	uc.UserUsecase = usecase.NewUserUsecase(
		lock.NewRedisLock(infras.Redis, "user-lock", 10*time.Second),
		repositories.UserRepository,
		domains.UserDomain,
	)

	return uc, nil
}
