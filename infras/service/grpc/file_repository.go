package grpc

import (
	"context"

	"github.com/todennus/proto/gen/service"
	"github.com/todennus/proto/gen/service/dto"
	"github.com/todennus/shared/authentication"
	"github.com/todennus/shared/enumdef"
	"github.com/todennus/shared/errordef"
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

func (repo *FileRepository) ValidateTemporaryFile(ctx context.Context, temporaryFileToken string) (snowflake.ID, error) {
	req := &dto.FileValidateTemporaryFileRequest{TemporaryFileToken: temporaryFileToken}
	resp, err := repo.fileClient.ValidateTemporaryFile(repo.auth.Context(ctx), req)
	if err != nil {
		return 0, errordef.ConvertGRPCError(err)
	}

	userID, err := snowflake.ParseString(resp.PolicyMetadata)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (repo *FileRepository) SaveToPersistent(ctx context.Context, temporaryFileToken string) (string, error) {
	req := &dto.FileCommandTemporaryFileRequest{
		Command:            enumdef.TemporaryFileCommandToGRPC(enumdef.TemporaryFileCommandSaveAsImage),
		TemporaryFileToken: temporaryFileToken,
		PolicySource:       enumdef.PolicySourceUserAvatar,
	}
	resp, err := repo.fileClient.CommandTemporaryFile(repo.auth.Context(ctx), req)
	if err != nil {
		return "", errordef.ConvertGRPCError(err)
	}

	return resp.PersistentUrl, nil
}

func (repo *FileRepository) DeleteTemporary(ctx context.Context, temporaryFileToken string) error {
	req := &dto.FileCommandTemporaryFileRequest{
		Command:            enumdef.TemporaryFileCommandToGRPC(enumdef.TemporaryFileCommandDelete),
		TemporaryFileToken: temporaryFileToken,
		PolicySource:       enumdef.PolicySourceUserAvatar,
	}
	_, err := repo.fileClient.CommandTemporaryFile(repo.auth.Context(ctx), req)
	if err != nil {
		return errordef.ConvertGRPCError(err)
	}

	return nil
}
