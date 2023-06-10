// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: ab5a6c03d7
// Version Date: 2023-06-04T17:09:20Z

// Package grpc provides a gRPC client for the User service.
package grpc

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	// This Service
	pb "github.com/POABOB/go-microservice/ticket-system/pb/user"
	"github.com/POABOB/go-microservice/ticket-system/user-service/svc"
)

// New returns an service backed by a gRPC client connection. It is the
// responsibility of the caller to dial, and later close, the connection.
func New(conn *grpc.ClientConn, options ...ClientOption) (pb.UserServer, error) {
	var cc clientConfig

	for _, f := range options {
		err := f(&cc)
		if err != nil {
			return nil, errors.Wrap(err, "cannot apply option")
		}
	}

	clientOptions := []grpctransport.ClientOption{
		grpctransport.ClientBefore(
			contextValuesToGRPCMetadata(cc.headers)),
	}
	var loginEndpoint endpoint.Endpoint
	{
		loginEndpoint = grpctransport.NewClient(
			conn,
			"user.User",
			"Login",
			EncodeGRPCLoginRequest,
			DecodeGRPCLoginResponse,
			pb.UserLoginResponse{},
			clientOptions...,
		).Endpoint()
	}

	var registerEndpoint endpoint.Endpoint
	{
		registerEndpoint = grpctransport.NewClient(
			conn,
			"user.User",
			"Register",
			EncodeGRPCRegisterRequest,
			DecodeGRPCRegisterResponse,
			pb.UserRegisterResponse{},
			clientOptions...,
		).Endpoint()
	}

	var loginwithgoogleEndpoint endpoint.Endpoint
	{
		loginwithgoogleEndpoint = grpctransport.NewClient(
			conn,
			"user.User",
			"LoginWithGoogle",
			EncodeGRPCLoginWithGoogleRequest,
			DecodeGRPCLoginWithGoogleResponse,
			pb.UserLoginResponse{},
			clientOptions...,
		).Endpoint()
	}

	var loginwithgooglecallbackEndpoint endpoint.Endpoint
	{
		loginwithgooglecallbackEndpoint = grpctransport.NewClient(
			conn,
			"user.User",
			"LoginWithGoogleCallback",
			EncodeGRPCLoginWithGoogleCallbackRequest,
			DecodeGRPCLoginWithGoogleCallbackResponse,
			pb.UserLoginResponse{},
			clientOptions...,
		).Endpoint()
	}

	var healthcheckEndpoint endpoint.Endpoint
	{
		healthcheckEndpoint = grpctransport.NewClient(
			conn,
			"user.User",
			"HealthCheck",
			EncodeGRPCHealthCheckRequest,
			DecodeGRPCHealthCheckResponse,
			pb.HealthCheckResponse{},
			clientOptions...,
		).Endpoint()
	}

	return svc.Endpoints{
		LoginEndpoint:                   loginEndpoint,
		RegisterEndpoint:                registerEndpoint,
		LoginWithGoogleEndpoint:         loginwithgoogleEndpoint,
		LoginWithGoogleCallbackEndpoint: loginwithgooglecallbackEndpoint,
		HealthCheckEndpoint:             healthcheckEndpoint,
	}, nil
}

// GRPC Client Decode

// DecodeGRPCLoginResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC login reply to a user-domain login response. Primarily useful in a client.
func DecodeGRPCLoginResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UserLoginResponse)
	return reply, nil
}

// DecodeGRPCRegisterResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC register reply to a user-domain register response. Primarily useful in a client.
func DecodeGRPCRegisterResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UserRegisterResponse)
	return reply, nil
}

// DecodeGRPCLoginWithGoogleResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC loginwithgoogle reply to a user-domain loginwithgoogle response. Primarily useful in a client.
func DecodeGRPCLoginWithGoogleResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UserLoginResponse)
	return reply, nil
}

// DecodeGRPCLoginWithGoogleCallbackResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC loginwithgooglecallback reply to a user-domain loginwithgooglecallback response. Primarily useful in a client.
func DecodeGRPCLoginWithGoogleCallbackResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UserLoginResponse)
	return reply, nil
}

// DecodeGRPCHealthCheckResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC healthcheck reply to a user-domain healthcheck response. Primarily useful in a client.
func DecodeGRPCHealthCheckResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.HealthCheckResponse)
	return reply, nil
}

// GRPC Client Encode

// EncodeGRPCLoginRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain login request to a gRPC login request. Primarily useful in a client.
func EncodeGRPCLoginRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UserLoginRequest)
	return req, nil
}

// EncodeGRPCRegisterRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain register request to a gRPC register request. Primarily useful in a client.
func EncodeGRPCRegisterRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UserRegisterRequest)
	return req, nil
}

// EncodeGRPCLoginWithGoogleRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain loginwithgoogle request to a gRPC loginwithgoogle request. Primarily useful in a client.
func EncodeGRPCLoginWithGoogleRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UserLoginRequest)
	return req, nil
}

// EncodeGRPCLoginWithGoogleCallbackRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain loginwithgooglecallback request to a gRPC loginwithgooglecallback request. Primarily useful in a client.
func EncodeGRPCLoginWithGoogleCallbackRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UserLoginRequest)
	return req, nil
}

// EncodeGRPCHealthCheckRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain healthcheck request to a gRPC healthcheck request. Primarily useful in a client.
func EncodeGRPCHealthCheckRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.HealthCheckRequest)
	return req, nil
}

type clientConfig struct {
	headers []string
}

// ClientOption is a function that modifies the client config
type ClientOption func(*clientConfig) error

func CtxValuesToSend(keys ...string) ClientOption {
	return func(o *clientConfig) error {
		o.headers = keys
		return nil
	}
}

func contextValuesToGRPCMetadata(keys []string) grpctransport.ClientRequestFunc {
	return func(ctx context.Context, md *metadata.MD) context.Context {
		var pairs []string
		for _, k := range keys {
			if v, ok := ctx.Value(k).(string); ok {
				pairs = append(pairs, k, v)
			}
		}

		if pairs != nil {
			*md = metadata.Join(*md, metadata.Pairs(pairs...))
		}

		return ctx
	}
}