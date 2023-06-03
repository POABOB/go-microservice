package service

import (
	"context"

	pb "github.com/POABOB/go-microservice/ticket-system/pb/user"
)

// pb.UserServer 已經有 interface 了
type UserService struct{}

func NewService() pb.UserServer {
	return UserService{}
}

func (s UserService) Login(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	// userEntity := model.NewUserModel()
	// res, err := userEntity.CheckUser(username, password)
	// if err != nil {
	// 	log.Printf("UserEntity.CreateUser, err : %v", err)
	// 	return 0, err
	// }
	// return res.UserId, nil

	return &pb.UserLoginResponse{}, nil
}

func (s UserService) Register(ctx context.Context, in *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	var resp pb.UserRegisterResponse
	return &resp, nil
}

func (s UserService) LoginWithGoogle(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	var resp pb.UserLoginResponse
	return &resp, nil
}

func (s UserService) LoginWithGoogleCallback(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	var resp pb.UserLoginResponse
	return &resp, nil
}

func (s UserService) HealthCheck(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	resp := &pb.HealthCheckResponse{Result: true, Err: ""}
	return resp, nil
}

// ServiceMiddleware define service middleware
type ServiceMiddleware func(pb.UserServer) pb.UserServer
