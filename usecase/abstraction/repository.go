package abstraction

import (
	"context"

	"github.com/todennus/shared/enumdef"
	"github.com/todennus/user-service/domain"
	"github.com/todennus/x/enum"
	"github.com/xybor-x/snowflake"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, userID snowflake.ID) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	CountByRole(ctx context.Context, role enum.Enum[enumdef.UserRole]) (int64, error)
}
