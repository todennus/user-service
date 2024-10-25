package grpc

import (
	"context"

	"github.com/todennus/proto/gen/service"
	pbdto "github.com/todennus/proto/gen/service/dto"
	"github.com/todennus/shared/errordef"
	"github.com/todennus/shared/response"
	"github.com/todennus/user-service/adapter/abstraction"
	"github.com/todennus/user-service/adapter/grpc/conversion"
	"google.golang.org/grpc/codes"
)

var _ service.UserServer = (*UserServer)(nil)

type UserServer struct {
	service.UnimplementedUserServer

	userUsecase abstraction.UserUsecase
}

func NewUserServer(userUsecase abstraction.UserUsecase) *UserServer {
	return &UserServer{
		userUsecase: userUsecase,
	}
}

func (s *UserServer) GetByID(ctx context.Context, req *pbdto.UserGetByIDRequest) (*pbdto.UserGetByIDResponse, error) {
	ucreq := conversion.NewUsecaseUserGetByIDRequest(req)
	resp, err := s.userUsecase.GetByID(ctx, ucreq)

	return response.NewGRPCResponseHandler(ctx, conversion.NewPbUserGetByIDResponse(resp), err).
		Map(codes.InvalidArgument, errordef.ErrRequestInvalid).
		Map(codes.NotFound, errordef.ErrNotFound).Finalize(ctx)
}

func (s *UserServer) Validate(ctx context.Context, req *pbdto.UserValidateRequest) (*pbdto.UserValidateResponse, error) {
	ucreq := conversion.NewUsecaseUserValidateRequest(req)
	resp, err := s.userUsecase.ValidateCredentials(ctx, ucreq)

	return response.NewGRPCResponseHandler(ctx, conversion.NewPbUserValidateResponse(resp), err).
		Map(codes.InvalidArgument, errordef.ErrRequestInvalid).
		Map(codes.PermissionDenied, errordef.ErrCredentialsInvalid).
		Map(codes.NotFound, errordef.ErrNotFound).Finalize(ctx)
}
