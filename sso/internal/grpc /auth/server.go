package auth

import (
	"context"

	"github.com/imhasandl/grpc-go/protos/gen/go/sso"
	"google.golang.org/grpc"
)

type serverAPI struct {
	sso.UnimplementedAuthServiceServer
}

func Register(gRPC *grpc.Server) {
	sso.RegisterAuthServiceServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, req *sso.LoginRequest) (*sso.LoginResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Register(ctx context.Context, req *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	panic("implement me")
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *sso.IsAdminRequest) (*sso.IsAdminResponse, error) {
	panic("implement me")
}
