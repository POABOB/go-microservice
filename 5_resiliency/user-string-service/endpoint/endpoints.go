package endpoint

import (
	"context"
	"errors"
	"use-string-service/service"

	"github.com/go-kit/kit/endpoint"
)

// 定義 StringEndpoints
type UseStringEndpoints struct {
	UseStringEndpoint   endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

// 請求錯誤
var (
	ErrInvalidRequestType = errors.New("RequestType has only two type: Concat, Diff")
)

// 定義請求結構體 StringRequest
type UseStringRequest struct {
	RequestType string `json:"request_type"`
	A           string `json:"a"`
	B           string `json:"b"`
}

// 定義回應結構體 StringResponse
type UseStringResponse struct {
	Result string `json:"result"`
	Error  error  `json:"error"`
}

// 服務請求的分發
func MakeUseStringEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UseStringRequest)

		var (
			res, a, b string
			opError   error
		)

		a = req.A
		b = req.B

		res, opError = svc.UseStringService(req.RequestType, a, b)

		return UseStringResponse{Result: res}, opError
	}
}

// 定義請求結構體 HealthRequest
type HealthRequest struct{}

// 定義回應結構體 HealthResponse
type HealthResponse struct {
	Status bool `json:"status"`
}

// 健康請求的分發
func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return HealthResponse{svc.HealthCheck()}, nil
	}
}
