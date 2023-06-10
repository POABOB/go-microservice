package handlers

import (
	"fmt"
	"time"

	pb "github.com/POABOB/go-microservice/ticket-system/pb/user"
	localconfig "github.com/POABOB/go-microservice/ticket-system/user-service/config"
	endpts "github.com/POABOB/go-microservice/ticket-system/user-service/endpoint"
	"github.com/POABOB/go-microservice/ticket-system/user-service/middlewares"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/time/rate"
)

// WrapEndpoints accepts the service's entire collection of endpoints, so that a
// set of middlewares can be wrapped around every middleware (e.g., access
// logging and instrumentation), and others wrapped selectively around some
// endpoints and not others (e.g., endpoints requiring authenticated access).
// Note that the final middleware wrapped will be the outermost middleware
// (i.e. applied first)
func WrapEndpoints(in endpts.Endpoints) endpts.Endpoints {

	// Pass a middleware you want applied to every endpoint.
	// optionally pass in endpoints by name that you want to be excluded
	// e.g.
	// in.WrapAllExcept(authMiddleware, "Status", "Ping")

	// Pass in a svc.LabeledMiddleware you want applied to every endpoint.
	// These middlewares get passed the endpoints name as their first argument when applied.
	// This can be used to write generic metric gathering middlewares that can
	// report the endpoint name for free.
	// github.com/metaverse/truss/_example/middlewares/labeledmiddlewares.go for examples.
	// in.WrapAllLabeledExcept(errorCounter(statsdCounter), "Status", "Ping")

	// How to apply a middleware to a single endpoint.
	// in.ExampleEndpoint = authMiddleware(in.ExampleEndpoint)

	// 設定 Ratelimiter，每秒最多 100 筆請求，
	ratebucket := rate.NewLimiter(rate.Every(time.Second*1), 500)
	in.WrapAllExcept(middlewares.TokenBucketLimitterWithBuildIn(ratebucket)) // 只有這個是外層

	// 設定 鏈路追蹤
	in.WrapAllLabeledExcept(middlewares.ZipkinMiddleware(localconfig.ZipkinTracer)) // 只有這個是外層

	// 設定 Prometheus
	fieldKeys := []string{"endpoint"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "kit_microservices",
		Subsystem: "user",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	in.WrapAllLabeledExcept(middlewares.Counter(requestCount))

	requestErrorCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "kit_microservices",
		Subsystem: "user",
		Name:      "request_error_count",
		Help:      "Number of errors occured.",
	}, fieldKeys)
	in.WrapAllLabeledExcept(middlewares.ErrorCounter(requestErrorCount))

	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "kit_microservices",
		Subsystem: "user",
		Name:      "request_latency",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	in.WrapAllLabeledExcept(middlewares.Latency(requestLatency))
	localconfig.Logger.Info(fmt.Sprintf("%v: %v, %v: %v", fieldKeys, requestCount, requestErrorCount, requestLatency))

	// 新增 Logging middleware
	in.WrapAllLabeledExcept(middlewares.Logging(localconfig.Logger))

	return in
}

func WrapService(in pb.UserServer) pb.UserServer {
	return in
}
