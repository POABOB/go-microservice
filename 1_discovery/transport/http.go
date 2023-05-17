package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"

	endpoints "ch6-discovery/endpoint"
)

var (
	ErrorBadRequest = errors.New("invalid reques parameter")
)

// 使用mux來實作http handler
func MakeHttpHandler(_ context.Context, endpoints endpoints.DiscoveryEndpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	// 定義錯誤處理器
	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	// SayHello路由
	r.Methods("GET").Path("/say-hello").Handler(kithttp.NewServer(
		endpoints.SayHelloEndpoint,
		decodeSayHelloRequest,
		encodeJsonResponse,
		options...,
	))

	// Discovery路由
	r.Methods("GET").Path("/discovery").Handler(kithttp.NewServer(
		endpoints.DiscoveryEndpoint,
		decodeDiscoveryRequest,
		encodeJsonResponse,
		options...,
	))

	// HealthCheck路由
	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeJsonResponse,
		options...,
	))

	return r
}

// 解碼Discovery的請求
func decodeDiscoveryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	serviceName := r.URL.Query().Get("serviceName")
	if serviceName == "" {
		return nil, ErrorBadRequest
	}
	return endpoints.DiscoveryRequest{ServiceName: serviceName}, nil
}

// 解碼SayHello的請求
func decodeSayHelloRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return endpoints.SayHelloRequest{}, nil
}

// 解碼HealthCheck的請求
func decodeHealthCheckRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return endpoints.HealthCheckRequest{}, nil
}

// 編碼Json Response
func encodeJsonResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// 編碼Error Response
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
