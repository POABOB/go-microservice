//go:build go1.7
// +build go1.7

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"zipkin-go/svc1"

	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"

	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	//"github.com/openzipkin-contrib/zipkin-go-opentracing/examples/string-services/svc1"
)

const (
	serviceName        string = "client"                             // 服務名稱
	hostPort           string = "0.0.0.0"                            // 地址+端口
	zipkinHTTPEndpoint string = "http://127.0.0.1:9411/api/v2/spans" // zipkin 的 endpoint
	debug              bool   = false                                // 是否 debug 開啟
	svc1Endpoint       string = "http://localhost:61001"             // 服務2的 endpoint
	sameSpan           bool   = true                                 // 相同的 span 可以設定成 RPC 風格的 spans (Zipkin V1) vs Node 風格 (OpenTracing)
	traceID128Bit      bool   = true                                 // 128 bits 的 trace ID 給根 span
)

// ci
func main() {
	// 設定收集的服務
	reporter := zipkinhttp.NewReporter(zipkinHTTPEndpoint)
	defer reporter.Close()

	endpoint, err := zipkin.NewEndpoint(serviceName, hostPort)
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
		os.Exit(-1)
	}

	sampler, err := zipkin.NewCountingSampler(1)
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
		os.Exit(-1)
	}

	// 設定配置，建立 tracer
	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithTraceID128Bit(traceID128Bit), zipkin.WithSampler(sampler))
	if err != nil {
		fmt.Printf("unable to create Zipkin tracer: %+v\n", err)
		os.Exit(-1)
	}

	// 使用 zipkin-go-opentracing 包覆我們的 tracer
	tracer := zipkinot.Wrap(nativeTracer)

	// 將我們的 tracer 設定為單例模式
	opentracing.SetGlobalTracer(tracer)

	// Create Client to svc1 Service
	client := svc1.NewHTTPClient(tracer, svc1Endpoint)

	// Create Root Span for duration of the interaction with svc1
	span := opentracing.StartSpan("Run")
	defer span.Finish()

	// Put root span in context so it will be used in our calls to the client.
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	// Call the Concat Method
	span.LogEvent("Call Concat")
	res1, err := client.Concat(ctx, "Hello", " World!")
	fmt.Printf("Concat: %s Err: %+v\n", res1, err)

	// Call the Sum Method
	span.LogEvent("Call Sum")
	res2, err := client.Sum(ctx, 10, 20)
	fmt.Printf("Sum: %d Err: %+v\n", res2, err)
}
