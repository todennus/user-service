package domain

import (
	"github.com/xybor-x/snowflake"
)

type AvatarPolicy struct {
	UserID       snowflake.ID
	AllowedTypes []string
	MaxSize      int64
}

type AvatarDomain struct {
	AllowedTypes []string
	MaxSize      int64
}

func NewAvatarDomain(
	allowedTypes []string,
	maxSize int64,
) *AvatarDomain {
	return &AvatarDomain{
		AllowedTypes: allowedTypes,
		MaxSize:      maxSize,
	}
}

func (domain *AvatarDomain) GetPolicy(userID snowflake.ID) *AvatarPolicy {
	return &AvatarPolicy{
		UserID:       userID,
		AllowedTypes: domain.AllowedTypes,
		MaxSize:      domain.MaxSize,
	}
}
