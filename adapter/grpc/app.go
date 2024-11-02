package grpc

import (
	"github.com/todennus/proto/gen/service"
	"github.com/todennus/shared/config"
	"github.com/todennus/shared/interceptor"
	"github.com/todennus/user-service/wiring"
	"google.golang.org/grpc"
)

func App(config *config.Config, usecases *wiring.Usecases) *grpc.Server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			interceptor.NewUnaryInterceptor().
				WithBasicContext().
				WithLogRoundTripTime().
				WithTimeout().
				WithAuthenticate().
				Interceptor(config),
		),
	)

	service.RegisterUserServer(s, NewUserServer(usecases.AvatarUsecase, usecases.UserUsecase))

	return s
}
