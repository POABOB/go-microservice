package string_service

import (
	"context"
	"errors"
	pb "go-kit/pb"
	"strings"

	"go-kit/config"

	"github.com/go-kit/kit/endpoint"
)

// CalculateEndpoint define endpoint
type StringEndpoints struct {
	StringEndpoint      endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

func (ue StringEndpoints) Concat(ctx context.Context, a string, b string) (string, error) {
	//ctx := context.Background()
	resp, err := ue.StringEndpoint(ctx, &pb.StringRequest{
		A: a,
		B: b,
	})

	response, ok := resp.(*pb.StringResponse)
	if !ok {
		config.Logger.Printf("%v", resp)
		return "", nil
	}
	return response.Ret, err
}

func (ue StringEndpoints) Diff(ctx context.Context, a string, b string) (string, error) {
	//ctx := context.Background()
	resp, err := ue.StringEndpoint(ctx, pb.StringRequest{
		A: a,
		B: b,
	})

	response, ok := resp.(*pb.StringResponse)
	if !ok {
		config.Logger.Printf("%v", resp)
		return "", nil
	}
	return response.Ret, err
}

// func (ue StringEndpoints) Check(ctx context.Context, serviceName string) bool {
// 	//ctx := context.Background()
// 	resp, _ := ue.HealthCheckEndpoint(ctx, &pb.HealthCheckRequest{
// 		Service: serviceName,
// 	})

// 	return resp.(*pb.HealthCheckResponse).Status
// }

var (
	ErrInvalidRequestType = errors.New("RequestType has only two type: Concat, Diff")
)

// StringRequest define request struct
type StringRequest struct {
	RequestType string `json:"request_type"`
	A           string `json:"a"`
	B           string `json:"b"`
}

// StringResponse define response struct
type StringResponse struct {
	Result string `json:"result"`
	Error  error  `json:"error"`
}

// MakeStringEndpoint make endpoint
func MakeStringEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(StringRequest)

		var (
			res, a, b string
			opError   error
		)

		a = req.A
		b = req.B

		if strings.EqualFold(req.RequestType, "Concat") {
			res, _ = svc.Concat(ctx, a, b)
		} else if strings.EqualFold(req.RequestType, "Diff") {
			res, _ = svc.Diff(ctx, a, b)
		} else {
			return nil, ErrInvalidRequestType
		}

		return StringResponse{Result: res, Error: opError}, nil
	}
}

// HealthRequest 健康检查请求结构
type HealthRequest struct{}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status bool `json:"status"`
}

// MakeHealthCheckEndpoint 创建健康检查Endpoint
func MakeHealthCheckEndpoint(svc Service) endpoint.Endpoint {
	return func(_ context.Context, _ interface{}) (response interface{}, err error) {
		status := true
		return HealthResponse{status}, nil
	}
}
