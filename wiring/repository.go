package wiring

import (
	"context"

	"github.com/todennus/user-service/infras/database/gorm"
	"github.com/todennus/user-service/infras/database/redis"
	"github.com/todennus/user-service/infras/service/grpc"
	"github.com/todennus/user-service/usecase/abstraction"
)

type Repositories struct {
	abstraction.UserRepository
	abstraction.FileRepository
	abstraction.AvatarPolicySessionRepository
}

func InitializeRepositories(ctx context.Context, infras *Infras) (*Repositories, error) {
	r := &Repositories{}

	r.UserRepository = gorm.NewUserRepository(infras.GormPostgres)
	r.FileRepository = grpc.NewFileRepository(infras.FilegRPCConn, infras.Auth)
	r.AvatarPolicySessionRepository = redis.NewAvatarPolicyRepository(infras.Redis)

	return r, nil
}
