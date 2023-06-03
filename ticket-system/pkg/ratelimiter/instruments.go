package main

import (
	"context"
	"errors"
	"time"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"
)

var ErrLimitExceed = errors.New("rate limit exceed")

// NewTokenBucketLimitterWithBuildIn 使用 x/time/rate 建立限流中間件
func NewTokenBucketLimitterWithBuildIn(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			// 超出限制
			if !bkt.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}

func DynamicLimitter(interval int, burst int) endpoint.Middleware {
	bucket := rate.NewLimiter(rate.Every(time.Second*time.Duration(interval)), burst)
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			// 超出限制
			if !bucket.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}
