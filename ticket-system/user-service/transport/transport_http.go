package transport

import (
	"encoding/json"
	"errors"
	"net/http"

	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	gozipkin "github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	// This service
	pb "github.com/POABOB/go-microservice/ticket-system/pb/user"
	"github.com/POABOB/go-microservice/ticket-system/user-service/endpoint"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

// MakeHttpHandler make http handler use mux
func MakeHttpHandler(_ context.Context, endpoints endpoint.UserEndpoints, zipkinTracer *gozipkin.Tracer, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	zipkinServer := zipkin.HTTPServerTrace(zipkinTracer, zipkin.Name("http-transport"))

	options := []kithttp.ServerOption{
		//kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		//kithttp.ServerErrorEncoder(kithttp.DefaultErrorEncoder),
		kithttp.ServerErrorEncoder(encodeError),
		zipkinServer,
	}

	r.Methods("POST").Path("/login").Handler(kithttp.NewServer(
		endpoints.LoginEndpoint,
		decodeLoginRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/register").Handler(kithttp.NewServer(
		endpoints.LoginEndpoint,
		decodeLoginRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/loginWithGoogle").Handler(kithttp.NewServer(
		endpoints.LoginEndpoint,
		decodeLoginRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/loginWithGoogleCallback").Handler(kithttp.NewServer(
		endpoints.LoginEndpoint,
		decodeLoginRequest,
		encodeResponse,
		options...,
	))

	r.Path("/metrics").Handler(promhttp.Handler())

	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeResponse,
		options...,
	))

	return r
}

// decodeLoginRequest decode request params to struct
func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req pb.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

// decodeRegisterRequest decode request params to struct
func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req pb.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

// decodeLoginWithGoogleRequest decode request params to struct
func decodeLoginWithGoogleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req pb.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

// decodeLoginWithGoogleCallbackRequest decode request params to struct
func decodeLoginWithGoogleCallbackRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req pb.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

// decodeHealthCheckRequest decode request params to struct
func decodeHealthCheckRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req pb.HealthCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

// encodeArithmeticResponse encode response to return
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
