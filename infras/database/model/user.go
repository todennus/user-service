package model

import (
	"time"

	"github.com/todennus/shared/enumdef"
	"github.com/todennus/user-service/domain"
	"github.com/todennus/x/enum"
	"github.com/xybor-x/snowflake"
)

type UserModel struct {
	ID          int64     `gorm:"id"`
	DisplayName string    `gorm:"display_name"`
	Username    string    `gorm:"username"`
	HashedPass  string    `gorm:"hashed_pass"`
	Role        string    `gorm:"role"`
	AvatarURL   string    `gorm:"avatar_url"`
	UpdatedAt   time.Time `gorm:"updated_at"`
}

func (UserModel) TableName() string {
	return "users"
}

func NewUser(d *domain.User) *UserModel {
	return &UserModel{
		ID:          d.ID.Int64(),
		DisplayName: d.DisplayName,
		Username:    d.Username,
		HashedPass:  d.HashedPass,
		AvatarURL:   d.AvatarURL,
		Role:        d.Role.String(),
		UpdatedAt:   d.UpdatedAt,
	}
}

func (u UserModel) To() (*domain.User, error) {
	return &domain.User{
		ID:          snowflake.ID(u.ID),
		DisplayName: u.DisplayName,
		Username:    u.Username,
		HashedPass:  u.HashedPass,
		Role:        enum.FromStr[enumdef.UserRole](u.Role),
		AvatarURL:   u.AvatarURL,
		UpdatedAt:   u.UpdatedAt,
	}, nil
}
