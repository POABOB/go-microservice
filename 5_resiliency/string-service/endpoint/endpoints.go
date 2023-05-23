package endpoint

import (
	"context"
	"errors"
	"string-service/service"
	"strings"

	"github.com/go-kit/kit/endpoint"
)

// 定義 StringEndpoints
type StringEndpoints struct {
	StringEndpoint      endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

// 請求錯誤
var (
	ErrInvalidRequestType = errors.New("RequestType has only two type: Concat, Diff")
)

// 定義請求結構體 StringRequest
type StringRequest struct {
	RequestType string `json:"request_type"`
	A           string `json:"a"`
	B           string `json:"b"`
}

// 定義回應結構體 StringResponse
type StringResponse struct {
	Result string `json:"result"`
	Error  error  `json:"error"`
}

// 服務請求的分發
func MakeStringEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(StringRequest)

		var (
			res, a, b string
			opError   error
		)

		a = req.A
		b = req.B
		// 根據RequestType判斷執行的函數
		if strings.EqualFold(req.RequestType, "Concat") {
			res, _ = svc.Concat(a, b)
		} else if strings.EqualFold(req.RequestType, "Diff") {
			res, _ = svc.Diff(a, b)
		} else {
			return nil, ErrInvalidRequestType
		}

		return StringResponse{Result: res, Error: opError}, nil
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
