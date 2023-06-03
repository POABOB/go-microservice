package endpoint

import (
	"time"

	pb "github.com/POABOB/go-microservice/ticket-system/pb/user"

	localconfig "github.com/POABOB/go-microservice/ticket-system/user-service/config"

	"github.com/POABOB/go-microservice/ticket-system/user-service/plugins"

	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"golang.org/x/time/rate"
)

func NewEndpoints(service pb.UserServer) UserEndpoints {
	// Business domain.

	// 設定 Ratelimiter，每秒最多 100 筆請求，
	ratebucket := rate.NewLimiter(rate.Every(time.Second*1), 100)

	// Endpoint domain.
	var (
		loginEndpoint = kitzipkin.TraceEndpoint(localconfig.ZipkinTracer, "login-endpoint")(
			plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(
				MakeLoginEndpoint(service),
			),
		)
		registerEndpoint = kitzipkin.TraceEndpoint(localconfig.ZipkinTracer, "register-endpoint")(
			plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(
				MakeRegisterEndpoint(service),
			),
		)
		loginwithgoogleEndpoint = kitzipkin.TraceEndpoint(localconfig.ZipkinTracer, "loginWithGoogle-endpoint")(
			plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(
				MakeLoginWithGoogleEndpoint(service),
			),
		)
		loginwithgooglecallbackEndpoint = kitzipkin.TraceEndpoint(localconfig.ZipkinTracer, "loginWithGoogleCallback-endpoint")(
			plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(
				MakeLoginWithGoogleCallbackEndpoint(service),
			),
		)
		healthcheckEndpoint = kitzipkin.TraceEndpoint(localconfig.ZipkinTracer, "healthCheck-endpoint")(
			plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(
				MakeHealthCheckEndpoint(service),
			),
		)
	)

	endpoints := UserEndpoints{
		LoginEndpoint:                   loginEndpoint,
		RegisterEndpoint:                registerEndpoint,
		LoginWithGoogleEndpoint:         loginwithgoogleEndpoint,
		LoginWithGoogleCallbackEndpoint: loginwithgooglecallbackEndpoint,
		HealthCheckEndpoint:             healthcheckEndpoint,
	}

	return endpoints
}
