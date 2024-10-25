package abstraction

import (
	"github.com/todennus/user-service/domain"
)

type UserDomain interface {
	Create(username, password string) (*domain.User, error)
	Validate(hashedPassword, password string) error
}
