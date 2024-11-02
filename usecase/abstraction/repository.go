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

	UpdateAvatarByID(ctx context.Context, userID snowflake.ID, avatarURL string) error

	CountByRole(ctx context.Context, role enum.Enum[enumdef.UserRole]) (int64, error)
}

type AvatarPolicySessionRepository interface {
	Store(ctx context.Context, policy *domain.AvatarPolicySession) error
	Load(ctx context.Context, policyToken string) (*domain.AvatarPolicySession, error)
	Delete(ctx context.Context, policyToken string) error
}

type FileRepository interface {
	ValidateTemporaryFile(ctx context.Context, temporaryFileToken string) (snowflake.ID, error)
	SaveToPersistent(ctx context.Context, temporaryFileToken string) (string, error)
	DeleteTemporary(ctx context.Context, temporaryFileToken string) error
}
