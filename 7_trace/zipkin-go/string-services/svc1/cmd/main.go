//go:build go1.7
// +build go1.7

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"zipkin-go/svc1"
	"zipkin-go/svc2"

	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"

	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	// "zipkin-go/svc1"
	// "zipkin-go/svc2"
)

const (
	serviceName        string = "svc1"                               // 服務名稱
	hostPort           string = "127.0.0.1:61001"                    // 地址+端口
	zipkinHTTPEndpoint string = "http://127.0.0.1:9411/api/v2/spans" // zipkin 的 endpoint
	debug              bool   = false                                // 是否 debug 開啟
	svc2Endpoint       string = "http://localhost:61002"             // 服務2的 endpoint
	sameSpan           bool   = true                                 // 相同的 span 可以設定成 RPC 風格的 spans (Zipkin V1) vs Node 風格 (OpenTracing)
	traceID128Bit      bool   = true                                 // 128 bits 的 trace ID 給根 span
)

// svc1
func main() {
	// 設定收集的服務
	reporter := zipkinhttp.NewReporter(zipkinHTTPEndpoint)
	defer reporter.Close()

	endpoint, err := zipkin.NewEndpoint(serviceName, hostPort)
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// 設定配置，建立 tracer
	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithTraceID128Bit(traceID128Bit))
	if err != nil {
		fmt.Printf("unable to create Zipkin tracer: %+v\n", err)
		os.Exit(-1)
	}

	// 使用 zipkin-go-opentracing 包覆我們的 tracer
	tracer := zipkinot.Wrap(nativeTracer)

	// 將我們的 tracer 設定為單例模式
	opentracing.SetGlobalTracer(tracer)

	// 建立 svc2 的 client，並調用服務
	service := svc1.NewService(svc2.NewHTTPClient(tracer, svc2Endpoint))
	// 建立 HTTP Handler
	handler := svc1.NewHTTPHandler(tracer, service)
	// 開啟
	fmt.Printf("Starting %s on %s\n", serviceName, hostPort)
	http.ListenAndServe(hostPort, handler)
}
