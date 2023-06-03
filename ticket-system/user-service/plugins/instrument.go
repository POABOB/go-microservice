package plugins

import (
	"context"
	"errors"

	"github.com/POABOB/go-microservice/ticket-system/user-service/service"

	"time"

	pb "github.com/POABOB/go-microservice/ticket-system/pb/user"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/juju/ratelimit"
	"golang.org/x/time/rate"
)

var ErrLimitExceed = errors.New("rate limit exceed")

// NewTokenBucketLimitterWithJuju 使用 juju/ratelimit 建立中間件
func NewTokenBucketLimitterWithJuju(bkt *ratelimit.Bucket) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if bkt.TakeAvailable(1) == 0 {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}

// NewTokenBucketLimitterWithBuildIn 使用x/time/rate 建立中間件
func NewTokenBucketLimitterWithBuildIn(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bkt.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}

// metricMiddleware 定義監控中間件，注入 Service
// 新增監控指標項目：requestCount 和 requestLatency
type metricMiddleware struct {
	Service        pb.UserServer
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

// Metrics 封裝監控方法
func Metrics(requestCount metrics.Counter, requestLatency metrics.Histogram) service.ServiceMiddleware {
	return func(next pb.UserServer) pb.UserServer {
		return metricMiddleware{
			next,
			requestCount,
			requestLatency}
	}
}

func (mw metricMiddleware) Login(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	defer func(beign time.Time) {
		lvs := []string{"method", "Login"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(beign).Seconds())
	}(time.Now())

	return mw.Service.Login(ctx, in)
}

func (mw metricMiddleware) Register(ctx context.Context, in *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	defer func(beign time.Time) {
		lvs := []string{"method", "Register"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(beign).Seconds())
	}(time.Now())

	return mw.Service.Register(ctx, in)
}

func (mw metricMiddleware) LoginWithGoogle(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	defer func(beign time.Time) {
		lvs := []string{"method", "LoginWithGoogle"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(beign).Seconds())
	}(time.Now())

	return mw.Service.LoginWithGoogle(ctx, in)
}

func (mw metricMiddleware) LoginWithGoogleCallback(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	defer func(beign time.Time) {
		lvs := []string{"method", "LoginWithGoogleCallback"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(beign).Seconds())
	}(time.Now())

	return mw.Service.LoginWithGoogleCallback(ctx, in)
}

func (mw metricMiddleware) HealthCheck(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.Service.HealthCheck(ctx, in)
}
