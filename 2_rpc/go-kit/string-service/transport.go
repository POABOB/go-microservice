package string_service

import (
	"context"
	"errors"
	pb "go-kit/pb"

	"github.com/go-kit/kit/transport/grpc"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

type grpcServer struct {
	concat grpc.Handler
	diff   grpc.Handler
	// check  grpc.Handler
}

// 只做資料傳遞
func (s *grpcServer) Concat(ctx context.Context, r *pb.StringRequest) (*pb.StringResponse, error) {
	_, resp, err := s.concat.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.StringResponse), nil
}

// 只做資料傳遞
func (s *grpcServer) Diff(ctx context.Context, r *pb.StringRequest) (*pb.StringResponse, error) {
	_, resp, err := s.diff.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.StringResponse), nil
}

// // 只做資料傳遞
// func (s *grpcServer) Check(ctx context.Context, r *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
// 	_, resp, err := s.check.ServeGRPC(ctx, r)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return resp.(*pb.HealthCheckResponse), nil
// }

func NewStringServer(ctx context.Context, endpoints StringEndpoints) pb.StringServiceServer {
	return &grpcServer{
		concat: grpc.NewServer(
			endpoints.StringEndpoint,
			DecodeConcatStringRequest,
			EncodeStringResponse,
		),
		diff: grpc.NewServer(
			endpoints.StringEndpoint,
			DecodeDiffStringRequest,
			EncodeStringResponse,
		),
		// check: grpc.NewServer(
		// 	endpoints.HealthCheckEndpoint,
		// 	DecodeHealthCheckRequest,
		// 	EncodeHealthCheckResponse,
		// ),
	}
}

func DecodeConcatStringRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.StringRequest)
	return StringRequest{
		RequestType: "Concat",
		A:           string(req.A),
		B:           string(req.B),
	}, nil
}

func DecodeDiffStringRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.StringRequest)
	return StringRequest{
		RequestType: "Diff",
		A:           string(req.A),
		B:           string(req.B),
	}, nil
}

func EncodeStringResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(StringResponse)

	if resp.Error != nil {
		return &pb.StringResponse{
			Ret: resp.Result,
			Err: resp.Error.Error(),
		}, nil
	}

	return &pb.StringResponse{
		Ret: resp.Result,
		Err: "",
	}, nil
}

// func DecodeHealthCheckRequest(_ context.Context, r interface{}) (interface{}, error) {
// 	req := r.(*pb.HealthCheckRequest)
// 	return &pb.HealthCheckRequest{
// 		Service: req.Service,
// 	}, nil
// }

// func EncodeHealthCheckResponse(_ context.Context, r interface{}) (interface{}, error) {
// 	resp := r.(*pb.HealthCheckResponse)

// 	return &pb.HealthCheckResponse{
// 		Status: resp.Status,
// 	}, nil
// }
