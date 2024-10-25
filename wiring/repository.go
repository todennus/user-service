package wiring

import (
	"context"

	"github.com/todennus/user-service/infras/database/gorm"
	"github.com/todennus/user-service/usecase/abstraction"
)

type Repositories struct {
	abstraction.UserRepository
}

func InitializeRepositories(ctx context.Context, infras *Infras) (*Repositories, error) {
	r := &Repositories{}

	r.UserRepository = gorm.NewUserRepository(infras.GormPostgres)

	return r, nil
}
