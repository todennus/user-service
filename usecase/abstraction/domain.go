package abstraction

import (
	"github.com/todennus/user-service/domain"
)

type UserDomain interface {
	New(username, password string) (*domain.User, error)
	NewFirst(username, password string) (*domain.User, error)
	Validate(hashedPassword, password string) error
}
