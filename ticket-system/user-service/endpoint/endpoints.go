// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: c27db7f217
// Version Date: 2023-06-10T03:54:57Z

package endpoint

// This file contains methods to make individual endpoints from services,
// request and response types to serve those endpoints, as well as encoders and
// decoders for those types, for all of our supported transport serialization
// formats.

import (
	"context"
	"fmt"

	pb "github.com/POABOB/go-microservice/ticket-system/pb/user"

	kitendpoint "github.com/go-kit/kit/endpoint"
)

// Endpoints collects all of the endpoints that compose an add service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
//
// In a server, it's useful for functions that need to operate on a per-endpoint
// basis. For example, you might pass an Endpoints to a function that produces
// an http.Handler, with each method (endpoint) wired up to a specific path. (It
// is probably a mistake in design to invoke the Service methods on the
// Endpoints struct in a server.)
//
// In a client, it's useful to collect individually constructed endpoints into a
// single type that implements the Service interface. For example, you might
// construct individual endpoints using transport/http.NewClient, combine them into an Endpoints, and return it to the caller as a Service.
type Endpoints struct {
	LoginEndpoint                   kitendpoint.Endpoint
	RegisterEndpoint                kitendpoint.Endpoint
	LoginWithGoogleEndpoint         kitendpoint.Endpoint
	LoginWithGoogleCallbackEndpoint kitendpoint.Endpoint
	HealthCheckEndpoint             kitendpoint.Endpoint
}

// Endpoints

func (e Endpoints) Login(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	response, err := e.LoginEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.UserLoginResponse), nil
}

func (e Endpoints) Register(ctx context.Context, in *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	response, err := e.RegisterEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.UserRegisterResponse), nil
}

func (e Endpoints) LoginWithGoogle(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	response, err := e.LoginWithGoogleEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.UserLoginResponse), nil
}

func (e Endpoints) LoginWithGoogleCallback(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	response, err := e.LoginWithGoogleCallbackEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.UserLoginResponse), nil
}

func (e Endpoints) HealthCheck(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	response, err := e.HealthCheckEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.HealthCheckResponse), nil
}

// Make Endpoints

func MakeLoginEndpoint(s pb.UserServer) kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.UserLoginRequest)
		v, err := s.Login(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeRegisterEndpoint(s pb.UserServer) kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.UserRegisterRequest)
		v, err := s.Register(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeLoginWithGoogleEndpoint(s pb.UserServer) kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.UserLoginRequest)
		v, err := s.LoginWithGoogle(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeLoginWithGoogleCallbackEndpoint(s pb.UserServer) kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.UserLoginRequest)
		v, err := s.LoginWithGoogleCallback(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeHealthCheckEndpoint(s pb.UserServer) kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.HealthCheckRequest)
		v, err := s.HealthCheck(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// WrapAllExcept wraps each Endpoint field of struct Endpoints with a
// go-kit/kit/endpoint.Middleware.
// Use this for applying a set of middlewares to every endpoint in the service.
// Optionally, endpoints can be passed in by name to be excluded from being wrapped.
// WrapAllExcept(middleware, "Status", "Ping")
func (e *Endpoints) WrapAllExcept(middleware kitendpoint.Middleware, excluded ...string) {
	included := map[string]struct{}{
		"Login":                   {},
		"Register":                {},
		"LoginWithGoogle":         {},
		"LoginWithGoogleCallback": {},
		"HealthCheck":             {},
	}

	for _, ex := range excluded {
		if _, ok := included[ex]; !ok {
			panic(fmt.Sprintf("Excluded endpoint '%s' does not exist; see middlewares/endpoints.go", ex))
		}
		delete(included, ex)
	}

	for inc := range included {
		if inc == "Login" {
			e.LoginEndpoint = middleware(e.LoginEndpoint)
		}
		if inc == "Register" {
			e.RegisterEndpoint = middleware(e.RegisterEndpoint)
		}
		if inc == "LoginWithGoogle" {
			e.LoginWithGoogleEndpoint = middleware(e.LoginWithGoogleEndpoint)
		}
		if inc == "LoginWithGoogleCallback" {
			e.LoginWithGoogleCallbackEndpoint = middleware(e.LoginWithGoogleCallbackEndpoint)
		}
		if inc == "HealthCheck" {
			e.HealthCheckEndpoint = middleware(e.HealthCheckEndpoint)
		}
	}
}

// LabeledMiddleware will get passed the endpoint name when passed to
// WrapAllLabeledExcept, this can be used to write a generic metrics
// middleware which can send the endpoint name to the metrics collector.
type LabeledMiddleware func(string, kitendpoint.Endpoint) kitendpoint.Endpoint

// WrapAllLabeledExcept wraps each Endpoint field of struct Endpoints with a
// LabeledMiddleware, which will receive the name of the endpoint. See
// LabeldMiddleware. See method WrapAllExept for details on excluded
// functionality.
func (e *Endpoints) WrapAllLabeledExcept(middleware func(string, kitendpoint.Endpoint) kitendpoint.Endpoint, excluded ...string) {
	included := map[string]struct{}{
		"Login":                   {},
		"Register":                {},
		"LoginWithGoogle":         {},
		"LoginWithGoogleCallback": {},
		"HealthCheck":             {},
	}

	for _, ex := range excluded {
		if _, ok := included[ex]; !ok {
			panic(fmt.Sprintf("Excluded endpoint '%s' does not exist; see middlewares/endpoints.go", ex))
		}
		delete(included, ex)
	}

	for inc := range included {
		if inc == "Login" {
			e.LoginEndpoint = middleware("Login", e.LoginEndpoint)
		}
		if inc == "Register" {
			e.RegisterEndpoint = middleware("Register", e.RegisterEndpoint)
		}
		if inc == "LoginWithGoogle" {
			e.LoginWithGoogleEndpoint = middleware("LoginWithGoogle", e.LoginWithGoogleEndpoint)
		}
		if inc == "LoginWithGoogleCallback" {
			e.LoginWithGoogleCallbackEndpoint = middleware("LoginWithGoogleCallback", e.LoginWithGoogleCallbackEndpoint)
		}
		if inc == "HealthCheck" {
			e.HealthCheckEndpoint = middleware("HealthCheck", e.HealthCheckEndpoint)
		}
	}
}

func NewEndpoints(service pb.UserServer) Endpoints {
	// Endpoint domain.
	var (
		loginEndpoint                   = MakeLoginEndpoint(service)
		registerEndpoint                = MakeRegisterEndpoint(service)
		loginwithgoogleEndpoint         = MakeLoginWithGoogleEndpoint(service)
		loginwithgooglecallbackEndpoint = MakeLoginWithGoogleCallbackEndpoint(service)
		healthcheckEndpoint             = MakeHealthCheckEndpoint(service)
	)

	endpoints := Endpoints{
		LoginEndpoint:                   loginEndpoint,
		RegisterEndpoint:                registerEndpoint,
		LoginWithGoogleEndpoint:         loginwithgoogleEndpoint,
		LoginWithGoogleCallbackEndpoint: loginwithgooglecallbackEndpoint,
		HealthCheckEndpoint:             healthcheckEndpoint,
	}

	return endpoints
}
