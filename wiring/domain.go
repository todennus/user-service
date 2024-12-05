package wiring

import (
	"context"

	"github.com/todennus/shared/config"
	"github.com/todennus/user-service/domain"
	"github.com/todennus/user-service/usecase/abstraction"
)

type Domains struct {
	abstraction.UserDomain
	abstraction.AvatarDomain
}

func InitializeDomains(ctx context.Context, config *config.Config) (*Domains, error) {
	var err error
	domains := &Domains{}

	domains.UserDomain, err = domain.NewUserDomain(config.SnowflakeNode)
	if err != nil {
		return nil, err
	}

	domains.AvatarDomain = domain.NewAvatarDomain(
		config.Variable.User.AvatarAllowedTypes,
		config.Variable.User.AvatarMaxSize,
	)

	return domains, nil
}
