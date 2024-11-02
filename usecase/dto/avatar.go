package dto

import (
	"github.com/todennus/user-service/domain"
	"github.com/xybor-x/snowflake"
)

type AvatarGetPolicyTokenRequest struct{}

type AvatarGetPolicyTokenResponse struct {
	PolicyToken string
}

func NewAvatarGetPolicyTokenResponse(policyToken string) *AvatarGetPolicyTokenResponse {
	return &AvatarGetPolicyTokenResponse{PolicyToken: policyToken}
}

type AvatarValidatePolicyTokenRequest struct {
	PolicyToken string
}

type AvatarValidatePolicyTokenResponse struct {
	UserID       snowflake.ID
	AllowedTypes []string
	MaxSize      int
}

func NewAvatarValidatePolicyTokenResponse(policy *domain.AvatarPolicySession) *AvatarValidatePolicyTokenResponse {
	return &AvatarValidatePolicyTokenResponse{
		UserID:       policy.UserID,
		AllowedTypes: policy.AllowedTypes,
		MaxSize:      policy.MaxSize,
	}
}

type AvatarUpdateRequest struct {
	TemporaryFileToken string
}

type AvatarUpdateResponse struct {
	AvatarURL string
}

func NewAvatarUpdateResponse(avatarURL string) *AvatarUpdateResponse {
	return &AvatarUpdateResponse{
		AvatarURL: avatarURL,
	}
}
