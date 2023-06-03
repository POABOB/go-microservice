package handlers

import (
	"context"

	pb "github.com/POABOB/go-microservice/ticket-system/pb/user"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.UserServer {
	return userService{}
}

type userService struct{}

func (s userService) Login(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	var resp pb.UserLoginResponse
	return &resp, nil
}

func (s userService) Register(ctx context.Context, in *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	var resp pb.UserRegisterResponse
	return &resp, nil
}

func (s userService) LoginWithGoogle(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	var resp pb.UserLoginResponse
	return &resp, nil
}

func (s userService) LoginWithGoogleCallback(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	var resp pb.UserLoginResponse
	return &resp, nil
}

func (s userService) HealthCheck(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	var resp pb.HealthCheckResponse
	return &resp, nil
}
