package auth

import (
	"context"
	ssov1 "github.com/rustam-ahmadov/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverApi struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) { //just registers handler
	ssov1.RegisterAuthServer(gRPC, &serverApi{auth: auth})
}

const emptyValue = 0

func (s *serverApi) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		//todo
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverApi) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		//todo
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverApi) IsAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	if req.GetUserId() == emptyValue {
		return nil, status.Errorf(codes.InvalidArgument, "userID is required")
	}
	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		//todo
		return nil, status.Errorf(codes.Internal, "internal err")
	}
	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

// region private
func validateLogin(req *ssov1.LoginRequest) error {
	if err := validateEmailPassword(req.GetEmail(), req.GetPassword()); err != nil {
		return err
	}
	if req.AppId == emptyValue {
		return status.Errorf(codes.InvalidArgument, "app_id is required")
	}
	return nil
}

func validateRegister(req *ssov1.RegisterRequest) error {
	if err := validateEmailPassword(req.GetEmail(), req.GetPassword()); err != nil {
		return err
	}
	return nil
}

func validateEmailPassword(email, password string) error {
	if email == "" {
		return status.Errorf(codes.InvalidArgument, "email is required")
	}
	if password == "" {
		return status.Errorf(codes.InvalidArgument, "password is required")
	}
	return nil
}

//endregion
