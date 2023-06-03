package handlers

import (
	pb "github.com/POABOB/go-microservice/ticket-system/pb/user"
	"github.com/POABOB/go-microservice/ticket-system/user-service/endpoint"

	localconfig "github.com/POABOB/go-microservice/ticket-system/user-service/config"

	"github.com/POABOB/go-microservice/ticket-system/user-service/plugins"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

// WrapEndpoints accepts the service's entire collection of endpoints, so that a
// set of middlewares can be wrapped around every middleware (e.g., access
// logging and instrumentation), and others wrapped selectively around some
// endpoints and not others (e.g., endpoints requiring authenticated access).
// Note that the final middleware wrapped will be the outermost middleware
// (i.e. applied first)
func WrapEndpoints(in endpoint.UserEndpoints) endpoint.UserEndpoints {

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

	return in
}

func WrapService(in pb.UserServer) pb.UserServer {

	// 設定 Prometheus
	fieldKeys := []string{"method"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "github.com/POABOB/go-microservice/ticket-system",
		Subsystem: "user_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "github.com/POABOB/go-microservice/ticket-system",
		Subsystem: "user_service",
		Name:      "request_latency",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	// 新增 Logging middleware
	in = plugins.LoggingMiddleware(localconfig.Logger)(in)
	// 新增 Prometheus middleware
	in = plugins.Metrics(requestCount, requestLatency)(in)

	return in
}
