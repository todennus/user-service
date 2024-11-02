package model

import (
	"time"

	"github.com/todennus/user-service/domain"
	"github.com/xybor-x/snowflake"
)

type AvatarPolicyRecord struct {
	UserID       int64    `json:"uid"`
	AllowedTypes []string `json:"ats"`
	MaxSize      int      `json:"mxs"`
	ExpiresAt    int64    `json:"exp"`
}

func NewAvatarPolicyRecord(policy *domain.AvatarPolicySession) *AvatarPolicyRecord {
	return &AvatarPolicyRecord{
		UserID:       policy.UserID.Int64(),
		AllowedTypes: policy.AllowedTypes,
		MaxSize:      policy.MaxSize,
		ExpiresAt:    policy.ExpiresAt.Unix(),
	}
}

func (record *AvatarPolicyRecord) To(policyToken string) *domain.AvatarPolicySession {
	return &domain.AvatarPolicySession{
		PolicyToken:  policyToken,
		AllowedTypes: record.AllowedTypes,
		UserID:       snowflake.ID(record.UserID),
		MaxSize:      record.MaxSize,
		ExpiresAt:    time.Unix(record.ExpiresAt, 0),
	}
}
