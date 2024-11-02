package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/todennus/shared/errordef"
	"github.com/todennus/user-service/domain"
	"github.com/todennus/user-service/infras/database/model"
)

func avatarPolicyKey(token string) string {
	return fmt.Sprintf("users:avatar_policy:%s", token)
}

type AvatarPolicyRepository struct {
	redis *redis.Client
}

func NewAvatarPolicyRepository(redis *redis.Client) *AvatarPolicyRepository {
	return &AvatarPolicyRepository{redis: redis}
}

func (repo *AvatarPolicyRepository) Store(ctx context.Context, policy *domain.AvatarPolicySession) error {
	record := model.NewAvatarPolicyRecord(policy)
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return err
	}

	return errordef.ConvertRedisError(
		repo.redis.SetEx(
			ctx,
			avatarPolicyKey(policy.PolicyToken),
			recordJSON,
			time.Until(policy.ExpiresAt),
		).Err(),
	)
}

func (repo *AvatarPolicyRepository) Load(ctx context.Context, policyToken string) (*domain.AvatarPolicySession, error) {
	recordJSON, err := repo.redis.Get(ctx, avatarPolicyKey(policyToken)).Result()
	if err != nil {
		return nil, errordef.ConvertRedisError(err)
	}

	record := model.AvatarPolicyRecord{}
	if err := json.Unmarshal([]byte(recordJSON), &record); err != nil {
		return nil, err
	}

	return record.To(policyToken), nil
}
func (repo *AvatarPolicyRepository) Delete(ctx context.Context, policyToken string) error {
	return errordef.ConvertRedisError(repo.redis.Del(ctx, avatarPolicyKey(policyToken)).Err())
}
