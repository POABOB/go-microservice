package service

import (
	"context"
	"fmt"

	pb "github.com/POABOB/go-microservice/ticket-system/pb/user"
	localconfig "github.com/POABOB/go-microservice/ticket-system/user-service/config"
	"github.com/gookit/validate"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.UserServer {
	return UserService{}
}

type UserService struct{}

func (s UserService) Login(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	v := validate.Struct(in)
	if !v.Validate() {
		localconfig.Logger.Error(fmt.Sprint("%v", v.Errors))
		return nil, v.Errors.OneError()
	}

	var resp pb.UserLoginResponse
	// do something ...

	return &resp, nil
}

func (s UserService) Register(ctx context.Context, in *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	v := validate.Struct(in)
	if !v.Validate() {
		localconfig.Logger.Error(fmt.Sprint("%v", v.Errors))
		return nil, v.Errors.OneError()
	}

	var resp pb.UserRegisterResponse
	// do something ...

	return &resp, nil
}

func (s UserService) LoginWithGoogle(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	v := validate.Struct(in)
	if !v.Validate() {
		localconfig.Logger.Error(fmt.Sprint("%v", v.Errors))
		return nil, v.Errors.OneError()
	}

	var resp pb.UserLoginResponse
	// do something ...

	return &resp, nil
}

func (s UserService) LoginWithGoogleCallback(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	v := validate.Struct(in)
	if !v.Validate() {
		localconfig.Logger.Error(fmt.Sprint("%v", v.Errors))
		return nil, v.Errors.OneError()
	}

	var resp pb.UserLoginResponse
	// do something ...

	return &resp, nil
}

func (s UserService) HealthCheck(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	v := validate.Struct(in)
	if !v.Validate() {
		localconfig.Logger.Error(fmt.Sprint("%v", v.Errors))
		return nil, v.Errors.OneError()
	}

	var resp pb.HealthCheckResponse
	// do something ...

	return &resp, nil
}
