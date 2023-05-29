package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"zipkin-kit/string-service/endpoint"
	"zipkin-kit/string-service/service"
	"zipkin-kit/string-service/transport"

	"github.com/POABOB/go-microservice/common/config"
	"github.com/POABOB/go-microservice/common/discover"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	uuid "github.com/satori/go.uuid"
)

func main() {
	var err error
	// 命令行讀取參數，沒有就使用預設值
	var (
		// 服務資訊
		serviceName = flag.String("service.name", "string-service", "service name")
		serviceHost = flag.String("service.host", "127.0.0.1", "service host")
		servicePort = flag.Int("service.port", 10086, "service port")
		// consul資訊
		consulHost = flag.String("consul.host", "127.0.0.1", "consul host")
		consulPort = flag.Int("consul.port", 8500, "consul port")
		// zipkin 資訊
		zipkinURL = flag.String("zipkin.url", "http://127.0.0.1:9411/api/v2/spans", "Zipkin server url")
		// grpcAddr  = flag.String("grpc", ":9008", "gRPC listen address.")
	)

	flag.Parse()
	ctx := context.Background()
	errChan := make(chan error)

	// Kit Logger
	var logger = config.KitLogger

	// 鏈路追中
	var zipkinTracer *zipkin.Tracer
	{
		var (
			hostPort                        = *serviceHost + ":" + strconv.Itoa(*servicePort)
			serviceName                     = *serviceName
			useNoopTracer bool              = (*zipkinURL == "")
			reporter      reporter.Reporter = zipkinhttp.NewReporter(*zipkinURL)
		)
		defer reporter.Close()
		zEP, _ := zipkin.NewEndpoint(serviceName, hostPort)
		// WithNoopTracer 如果沒有任何配置資訊，會自動使用預設值
		if zipkinTracer, err = zipkin.NewTracer(
			reporter, zipkin.WithLocalEndpoint(zEP), zipkin.WithNoopTracer(useNoopTracer),
		); err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}

		if !useNoopTracer {
			logger.Log("tracer", "Zipkin", "type", "Native", "URL", *zipkinURL)
		}
	}

	// 服務發現Client
	// discoveryClient, err := discover.NewMyDiscoverClient(*consulHost, *consulPort)
	discoveryClient, err := discover.NewKitHTTPDiscoverClient(*consulHost, *consulPort)

	// 獲取服務失敗
	if err != nil {
		config.Logger.Println("Get Consul Client failed")
		os.Exit(-1)
	}

	// 把服務注入至 Endpoint
	var svc service.Service = service.StringService{}
	endpts := endpoint.StringEndpoints{
		StringEndpoint:      kitzipkin.TraceEndpoint(zipkinTracer, "string-endpoint")(endpoint.MakeStringEndpoint(ctx, svc)), // String Service的 Endpoint
		HealthCheckEndpoint: kitzipkin.TraceEndpoint(zipkinTracer, "health-endpoint")(endpoint.MakeHealthCheckEndpoint(svc)), // 健檢的Endpoint
	}

	// 建立 http.Handler
	r := transport.MakeHttpHandler(ctx, endpts, zipkinTracer, config.KitLogger)
	instanceId := *serviceName + "-" + uuid.NewV4().String()

	// 啟動 http server
	go func() {
		config.Logger.Println("Http Server start at port:" + strconv.Itoa(*servicePort))
		// 啟動前註冊Service
		if !discoveryClient.Register(*serviceName, instanceId, *serviceHost, "/health", *servicePort, nil, config.Logger) {
			config.Logger.Printf("string-service for service %s failed.", *serviceName)
			// 註冊失敗，啟動失敗
			os.Exit(-1)
		}
		handler := r
		errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), handler)
	}()
	// //grpc server
	// go func() {
	// 	fmt.Println("grpc Server start at port" + *grpcAddr)
	// 	listener, err := net.Listen("tcp", *grpcAddr)
	// 	if err != nil {
	// 		errChan <- err
	// 		return
	// 	}
	// 	serverTracer := kitzipkin.GRPCServerTrace(zipkinTracer, kitzipkin.Name("string-grpc-transport"))

	// 	handler := NewGRPCServer(ctx, endpts, serverTracer)
	// 	gRPCServer := grpc.NewServer()
	// 	pb.RegisterStringServiceServer(gRPCServer, handler)
	// 	errChan <- gRPCServer.Serve(listener)
	// }()

	// 監聽syscall，如果 ctrl + c 被通知，關閉服務
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	// 服務註銷
	discoveryClient.DeRegister(instanceId, config.Logger)
	config.Logger.Println(error)
}
