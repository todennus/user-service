package grpc

import (
	"context"
	"time"

	"github.com/todennus/proto/gen/service"
	"github.com/todennus/proto/gen/service/dto"
	"github.com/todennus/shared/authentication"
	"github.com/todennus/shared/errordef"
	"github.com/todennus/user-service/domain"
	"github.com/xybor-x/snowflake"
	"google.golang.org/grpc"
)

type FileRepository struct {
	auth       *authentication.GrpcAuthorization
	fileClient service.FileClient
}

func NewFileRepository(client *grpc.ClientConn, auth *authentication.GrpcAuthorization) *FileRepository {
	return &FileRepository{
		fileClient: service.NewFileClient(client),
		auth:       auth,
	}
}

func (repo *FileRepository) RegisterUpload(ctx context.Context, policy *domain.AvatarPolicy) (string, error) {
	req := &dto.FileRegisterUploadRequest{
		UserId:       policy.UserID.Int64(),
		AllowedTypes: policy.AllowedTypes,
		MaxSize:      policy.MaxSize,
	}
	resp, err := repo.fileClient.RegisterUpload(repo.auth.Context(ctx), req)
	if err != nil {
		return "", errordef.ConvertGRPCError(err)
	}

	return resp.UploadToken, nil
}

func (repo *FileRepository) CreatePresignedURL(ctx context.Context, ownershipID snowflake.ID, expiration time.Duration) (string, error) {
	req := &dto.FileCreatePresignedURLRequest{
		OwnershipId: ownershipID.Int64(),
		Expiration:  int64(expiration / time.Second),
	}
	resp, err := repo.fileClient.CreatePresignedURL(repo.auth.Context(ctx), req)
	if err != nil {
		return "", errordef.ConvertGRPCError(err)
	}

	return resp.PresignedUrl, nil
}

func (repo *FileRepository) ChangeRefcount(ctx context.Context, inc, dec []snowflake.ID) error {
	incOwnershipID := []int64{}
	for i := range inc {
		incOwnershipID = append(incOwnershipID, inc[i].Int64())
	}

	decOwnershipID := []int64{}
	for i := range dec {
		decOwnershipID = append(decOwnershipID, dec[i].Int64())
	}

	req := &dto.FileChangeRefcountRequest{
		IncOwnershipId: incOwnershipID,
		DecOwnershipId: decOwnershipID,
	}

	if _, err := repo.fileClient.ChangeRefcount(repo.auth.Context(ctx), req); err != nil {
		return errordef.ConvertGRPCError(err)
	}

	return nil
}
