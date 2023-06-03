package plugins

import (
	"context"
	"time"

	pb "github.com/POABOB/go-microservice/ticket-system/pb/user"
	"github.com/POABOB/go-microservice/ticket-system/user-service/service"

	"github.com/go-kit/log"
)

// loggingMiddleware Make a new type
// that contains Service interface and logger instance
type loggingMiddleware struct {
	Service pb.UserServer
	logger  log.Logger
}

// LoggingMiddleware make logging middleware
func LoggingMiddleware(logger log.Logger) service.ServiceMiddleware {
	return func(next pb.UserServer) pb.UserServer {
		return loggingMiddleware{next, logger}
	}
}

func (mw loggingMiddleware) Login(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	res, err := mw.Service.Login(ctx, in)
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"function", "Login",
			"request", in,
			"response", res,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return res, err
}

func (mw loggingMiddleware) Register(ctx context.Context, in *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	res, err := mw.Service.Register(ctx, in)
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"function", "Register",
			"request", in,
			"response", res,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return res, err
}

func (mw loggingMiddleware) LoginWithGoogle(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	res, err := mw.Service.LoginWithGoogle(ctx, in)
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"function", "LoginWithGoogle",
			"request", in,
			"response", res,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return res, err
}

func (mw loggingMiddleware) LoginWithGoogleCallback(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	res, err := mw.Service.LoginWithGoogleCallback(ctx, in)
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"function", "LoginWithGoogleCallback",
			"request", in,
			"response", res,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return res, err
}

func (mw loggingMiddleware) HealthCheck(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	res, err := mw.Service.HealthCheck(ctx, in)
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"function", "HealthChcek",
			"request", in,
			"response", res,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return res, err
}
