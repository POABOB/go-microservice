package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	pb "github.com/POABOB/go-microservice/ticket-system/pb/user"
)

type UserEndpoints struct {
	LoginEndpoint                   endpoint.Endpoint
	RegisterEndpoint                endpoint.Endpoint
	LoginWithGoogleEndpoint         endpoint.Endpoint
	LoginWithGoogleCallbackEndpoint endpoint.Endpoint
	HealthCheckEndpoint             endpoint.Endpoint
}

// Endpoints
func (e UserEndpoints) Login(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	response, err := e.LoginEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.UserLoginResponse), nil
}

func (e UserEndpoints) Register(ctx context.Context, in *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	response, err := e.RegisterEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.UserRegisterResponse), nil
}

func (e UserEndpoints) LoginWithGoogle(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	response, err := e.LoginWithGoogleEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.UserLoginResponse), nil
}

func (e UserEndpoints) LoginWithGoogleCallback(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	response, err := e.LoginWithGoogleCallbackEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.UserLoginResponse), nil
}

func (e UserEndpoints) HealthCheck(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	response, err := e.HealthCheckEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.HealthCheckResponse), nil
}

// Make Endpoints
func MakeLoginEndpoint(s pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.UserLoginRequest)
		v, err := s.Login(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeRegisterEndpoint(s pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.UserRegisterRequest)
		v, err := s.Register(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeLoginWithGoogleEndpoint(s pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.UserLoginRequest)
		v, err := s.LoginWithGoogle(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeLoginWithGoogleCallbackEndpoint(s pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.UserLoginRequest)
		v, err := s.LoginWithGoogleCallback(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeHealthCheckEndpoint(s pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.HealthCheckRequest)
		v, err := s.HealthCheck(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}
