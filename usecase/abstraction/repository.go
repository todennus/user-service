package abstraction

import (
	"context"
	"time"

	"github.com/todennus/shared/enumdef"
	"github.com/todennus/user-service/domain"
	"github.com/xybor-x/snowflake"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error

	GetByID(ctx context.Context, userID snowflake.ID) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)

	GetAvatarByID(ctx context.Context, userID snowflake.ID) (snowflake.ID, error)
	UpdateAvatarByID(ctx context.Context, userID, ownershipID snowflake.ID) error

	CountByRole(ctx context.Context, role enumdef.UserRole) (int64, error)
}

type FileRepository interface {
	RegisterUpload(ctx context.Context, policy *domain.AvatarPolicy) (string, error)
	CreatePresignedURL(ctx context.Context, ownershipID snowflake.ID, expiration time.Duration) (string, error)
	ChangeRefcount(ctx context.Context, incOwnershipID, decOwnershipID []snowflake.ID) error
}
