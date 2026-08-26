package auth_grpc

import (
	"context"

	"github.com/emp2ty0/sso-protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		app_id int32,
	) (token string, err error)

	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
}

type ServerAPI struct {
	sso.UnimplementedAuthServer
	authService AuthService
}

func Register(gRPC *grpc.Server, authService AuthService) {
	sso.RegisterAuthServer(gRPC, &ServerAPI{authService: authService})
}

func (s *ServerAPI) Login(
	ctx context.Context,
	req *sso.LoginRequest,
) (*sso.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.authService.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error ")
	}

	return &sso.LoginResponse{
		Token: token,
	}, nil
}

func (s *ServerAPI) Register(
	ctx context.Context,
	req *sso.RegisterRequest,
) (*sso.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.authService.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error ")
	}

	return &sso.RegisterResponse{
		UserId: userID,
	}, nil
}

func validateLogin(req *sso.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required ")
	}
	if req.GetAppId() == 0 {
		return status.Error(codes.InvalidArgument, "app_id is required ")
	}

	return nil
}

func validateRegister(req *sso.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required ")
	}

	return nil
}
