package model

import (
	"time"

	"github.com/todennus/shared/enumdef"
	"github.com/todennus/user-service/domain"
	"github.com/todennus/x/enum"
	"github.com/xybor-x/snowflake"
)

type UserModel struct {
	ID          int64     `gorm:"column:id"`
	DisplayName string    `gorm:"column:display_name"`
	Username    string    `gorm:"column:username"`
	HashedPass  string    `gorm:"column:hashed_pass"`
	Role        string    `gorm:"column:role"`
	Avatar      int64     `gorm:"column:avatar"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
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
		Avatar:      d.Avatar.Int64(),
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
		Avatar:      snowflake.ParseInt64(u.Avatar),
		UpdatedAt:   u.UpdatedAt,
	}, nil
}
