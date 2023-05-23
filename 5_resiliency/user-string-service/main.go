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

	"use-string-service/config"
	"use-string-service/endpoint"
	loadbalance "use-string-service/load-balance"
	"use-string-service/plugins"
	"use-string-service/service"
	"use-string-service/transport"

	"github.com/POABOB/go-microservice/common/discover"
	"github.com/go-kit/kit/circuitbreaker"
	uuid "github.com/satori/go.uuid"
)

func main() {
	// 命令行讀取參數，沒有就使用預設值
	var (
		// 服務資訊
		serviceName = flag.String("service.name", "UseStringService", "service name")
		serviceHost = flag.String("service.host", "127.0.0.1", "service host")
		servicePort = flag.Int("service.port", 10086, "service port")
		// consul資訊
		consulHost = flag.String("consul.host", "127.0.0.1", "consul host")
		consulPort = flag.Int("consul.port", 8500, "consul port")
	)

	flag.Parse()
	ctx := context.Background()
	errChan := make(chan error)

	// 服務發現Client
	discoveryClient, err := discover.NewKitHTTPDiscoverClient(*consulHost, *consulPort)

	// 獲取服務失敗
	if err != nil {
		config.Logger.Println("Get Consul Client failed")
		os.Exit(-1)
	}

	// 將服務注入 LoggingMiddleware
	var svc service.Service = service.NewUseStringService(discoveryClient, &loadbalance.RandomLoadBalance{})
	svc = plugins.LoggingMiddleware(config.KitLogger)(svc)

	// 把服務注入至 Endpoint
	endpts := endpoint.UseStringEndpoints{
		// 設定 Hystrix Middleware 來使用熔斷
		UseStringEndpoint:   circuitbreaker.Hystrix(service.StringServiceCommandName)(endpoint.MakeUseStringEndpoint(svc)), // String Service的 Endpoint
		HealthCheckEndpoint: endpoint.MakeHealthCheckEndpoint(svc),                                                         // 健檢的Endpoint
	}

	// 建立 http.Handler
	r := transport.MakeHttpHandler(ctx, endpts, config.KitLogger)
	instanceId := *serviceName + "-" + uuid.NewV4().String()

	// 啟動 http server
	go func() {
		config.Logger.Println("Http Server start at port:" + strconv.Itoa(*servicePort))
		// 啟動前註冊Service
		if !discoveryClient.Register(*serviceName, instanceId, *serviceHost, "/health", *servicePort, nil, config.Logger) {
			config.Logger.Printf("use-string-service for service %s failed.", *serviceName)
			// 註冊失敗，啟動失敗
			os.Exit(-1)
		}
		handler := r
		errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), handler)
	}()

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
