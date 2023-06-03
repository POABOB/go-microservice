package transport

import (
	"context"

	"github.com/go-kit/kit/transport/grpc"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	// This Service
	pb "github.com/POABOB/go-microservice/ticket-system/pb/user"
	"github.com/POABOB/go-microservice/ticket-system/user-service/endpoint"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC UserServer.
func MakeGRPCServer(ctx context.Context, endpoints endpoint.UserEndpoints, serverTracer grpc.ServerOption) pb.UserServer {
	return &grpcServer{
		login: grpctransport.NewServer(
			endpoints.LoginEndpoint,
			DecodeGRPCLoginRequest,
			EncodeGRPCLoginResponse,
			serverTracer,
		),
		register: grpctransport.NewServer(
			endpoints.RegisterEndpoint,
			DecodeGRPCRegisterRequest,
			EncodeGRPCRegisterResponse,
			serverTracer,
		),
		loginwithgoogle: grpctransport.NewServer(
			endpoints.LoginWithGoogleEndpoint,
			DecodeGRPCLoginWithGoogleRequest,
			EncodeGRPCLoginWithGoogleResponse,
			serverTracer,
		),
		loginwithgooglecallback: grpctransport.NewServer(
			endpoints.LoginWithGoogleCallbackEndpoint,
			DecodeGRPCLoginWithGoogleCallbackRequest,
			EncodeGRPCLoginWithGoogleCallbackResponse,
			serverTracer,
		),
		healthcheck: grpctransport.NewServer(
			endpoints.HealthCheckEndpoint,
			DecodeGRPCHealthCheckRequest,
			EncodeGRPCHealthCheckResponse,
			serverTracer,
		),
	}
}

// grpcServer implements the UserServer interface
type grpcServer struct {
	login                   grpctransport.Handler
	register                grpctransport.Handler
	loginwithgoogle         grpctransport.Handler
	loginwithgooglecallback grpctransport.Handler
	healthcheck             grpctransport.Handler
}

// Methods for grpcServer to implement UserServer interface

func (s *grpcServer) Login(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	_, rep, err := s.login.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UserLoginResponse), nil
}

func (s *grpcServer) Register(ctx context.Context, req *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	_, rep, err := s.register.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UserRegisterResponse), nil
}

func (s *grpcServer) LoginWithGoogle(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	_, rep, err := s.loginwithgoogle.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UserLoginResponse), nil
}

func (s *grpcServer) LoginWithGoogleCallback(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	_, rep, err := s.loginwithgooglecallback.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UserLoginResponse), nil
}

func (s *grpcServer) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	_, rep, err := s.healthcheck.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.HealthCheckResponse), nil
}
