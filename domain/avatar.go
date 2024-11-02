package domain

import (
	"time"

	"github.com/todennus/shared/enumdef"
	"github.com/todennus/x/xcrypto"
	"github.com/xybor-x/snowflake"
)

type AvatarPolicySession struct {
	PolicyToken  string
	UserID       snowflake.ID
	AllowedTypes []string
	MaxSize      int
	ExpiresAt    time.Time
}

type AvatarDomain struct {
	AllowedTypes          []string
	MaxSize               int
	PolicyTokenExpiration time.Duration
}

func NewAvatarDomain(
	allowedTypes []string,
	maxSize int,
	policyTokenExpiration time.Duration,
) *AvatarDomain {
	return &AvatarDomain{
		AllowedTypes:          allowedTypes,
		MaxSize:               maxSize,
		PolicyTokenExpiration: policyTokenExpiration,
	}
}

func (domain *AvatarDomain) GetPolicy(userID snowflake.ID) *AvatarPolicySession {
	return &AvatarPolicySession{
		PolicyToken:  enumdef.FilePolicyToken(enumdef.PolicySourceUserAvatar, xcrypto.RandToken()),
		UserID:       userID,
		AllowedTypes: domain.AllowedTypes,
		MaxSize:      domain.MaxSize,
		ExpiresAt:    time.Now().Add(domain.PolicyTokenExpiration),
	}
}
