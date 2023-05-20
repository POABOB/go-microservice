package endpoint

import (
	"1_discovery/service"
	"context"

	"github.com/go-kit/kit/endpoint"
)

type DiscoveryEndpoints struct {
	HealthCheckEndpoint endpoint.Endpoint
	SayHelloEndpoint    endpoint.Endpoint
	DiscoveryEndpoint   endpoint.Endpoint
}

// SayHello請求結構體
type SayHelloRequest struct{}

// SayHello回應結構體
type SayHelloResponse struct {
	Message string `json:"message"`
}

// 創建SayHello Endpoint
func MakeSayHelloEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return SayHelloResponse{
			Message: svc.SayHello(),
		}, nil
	}
}

// Discovery請求結構體
type DiscoveryRequest struct {
	ServiceName string
}

// Discovery回應結構體
type DiscoveryResponse struct {
	Instances []interface{} `json:"instances"`
	Error     string        `json:"error"`
}

// 創建Discovery Endpoint
func MakeDiscoveryEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DiscoveryRequest)
		instances, err := svc.DiscoveryService(ctx, req.ServiceName)

		var errorString = ""
		if err != nil {
			errorString = err.Error()
		}

		return &DiscoveryResponse{
			Instances: instances,
			Error:     errorString,
		}, nil
	}
}

// HealthCheck請求結構體
type HealthCheckRequest struct{}

// HealthCheck回應結構體
type HealthCheckResponse struct {
	Status bool `json:"status"`
}

// 創建HealthCheck Endpoint
func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return HealthCheckResponse{
			Status: svc.HealthCheck(),
		}, nil
	}
}
