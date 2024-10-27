package wiring

import (
	"context"

	"github.com/todennus/shared/config"
	"github.com/todennus/user-service/domain"
	"github.com/todennus/user-service/usecase/abstraction"
)

type Domains struct {
	abstraction.UserDomain
}

func InitializeDomains(ctx context.Context, config *config.Config) (*Domains, error) {
	var err error
	domains := &Domains{}

	domains.UserDomain, err = domain.NewUserDomain(config.SnowflakeNode)
	if err != nil {
		return nil, err
	}

	return domains, nil
}
