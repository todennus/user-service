package abstraction

import (
	"github.com/todennus/user-service/domain"
	"github.com/xybor-x/snowflake"
)

type UserDomain interface {
	New(username, password string) (*domain.User, error)
	NewFirst(username, password string) (*domain.User, error)
	Validate(hashedPassword, password string) error
}

type AvatarDomain interface {
	GetPolicy(userID snowflake.ID) *domain.AvatarPolicy
}
