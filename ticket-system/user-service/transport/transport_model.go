package transport

import (
	"context"

	pb "github.com/POABOB/go-microservice/ticket-system/pb/user"
)

// Server Decode

// DecodeGRPCLoginRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC login request to a user-domain login request. Primarily useful in a server.
func DecodeGRPCLoginRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UserLoginRequest)
	return req, nil
}

// DecodeGRPCRegisterRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC register request to a user-domain register request. Primarily useful in a server.
func DecodeGRPCRegisterRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UserRegisterRequest)
	return req, nil
}

// DecodeGRPCLoginWithGoogleRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC loginwithgoogle request to a user-domain loginwithgoogle request. Primarily useful in a server.
func DecodeGRPCLoginWithGoogleRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UserLoginRequest)
	return req, nil
}

// DecodeGRPCLoginWithGoogleCallbackRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC loginwithgooglecallback request to a user-domain loginwithgooglecallback request. Primarily useful in a server.
func DecodeGRPCLoginWithGoogleCallbackRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UserLoginRequest)
	return req, nil
}

// DecodeGRPCHealthCheckRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC healthcheck request to a user-domain healthcheck request. Primarily useful in a server.
func DecodeGRPCHealthCheckRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.HealthCheckRequest)
	return req, nil
}

// Server Encode

// EncodeGRPCLoginResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain login response to a gRPC login reply. Primarily useful in a server.
func EncodeGRPCLoginResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.UserLoginResponse)
	return resp, nil
}

// EncodeGRPCRegisterResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain register response to a gRPC register reply. Primarily useful in a server.
func EncodeGRPCRegisterResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.UserRegisterResponse)
	return resp, nil
}

// EncodeGRPCLoginWithGoogleResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain loginwithgoogle response to a gRPC loginwithgoogle reply. Primarily useful in a server.
func EncodeGRPCLoginWithGoogleResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.UserLoginResponse)
	return resp, nil
}

// EncodeGRPCLoginWithGoogleCallbackResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain loginwithgooglecallback response to a gRPC loginwithgooglecallback reply. Primarily useful in a server.
func EncodeGRPCLoginWithGoogleCallbackResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.UserLoginResponse)
	return resp, nil
}

// EncodeGRPCHealthCheckResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain healthcheck response to a gRPC healthcheck reply. Primarily useful in a server.
func EncodeGRPCHealthCheckResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.HealthCheckResponse)
	return resp, nil
}
