package resource

import (
	"context"

	"github.com/todennus/shared/scopedef"
	"github.com/todennus/shared/xcontext"
	"github.com/todennus/user-service/domain"
	"github.com/todennus/x/conversion"
	"github.com/xybor-x/snowflake"
)

type User struct {
	ID          snowflake.ID
	Username    *string
	DisplayName *string
	Role        *string
	AvatarURL   *string
}

func NewUserWithFilter(ctx context.Context, user *domain.User, avatarURL string) *User {
	usecaseUser := NewUser(user, avatarURL)

	scopedef.Eval(xcontext.Scope(ctx)).
		RequireAdmin(scopedef.AdminReadUserProfile).
		RequireUser(ctx, scopedef.UserReadUserProfile, user.ID).
		FilterIfUnsatisfied(&usecaseUser.Username, &usecaseUser.DisplayName)

	scopedef.Eval(xcontext.Scope(ctx)).
		RequireAdmin(scopedef.AdminReadUserProfile).
		FilterIfUnsatisfied(&usecaseUser.Role)

	scopedef.Eval(xcontext.Scope(ctx)).
		RequireAdmin(scopedef.AdminReadUserProfile).
		RequireUser(ctx, scopedef.UserReadUserAvatar, user.ID).
		FilterIfUnsatisfied(&usecaseUser.AvatarURL)

	return usecaseUser
}

func NewUser(user *domain.User, avatarURL string) *User {
	usecaseUser := &User{
		ID:          user.ID,
		Username:    &user.Username,
		DisplayName: &user.DisplayName,
		Role:        conversion.ConvertToPointer(user.Role.String()),
	}

	if avatarURL != "" {
		usecaseUser.AvatarURL = &avatarURL
	}

	return usecaseUser
}
