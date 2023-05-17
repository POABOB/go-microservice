package main

import (
	"ch6-discovery/config"
	"ch6-discovery/discover"
	"ch6-discovery/endpoint"
	"ch6-discovery/service"
	"ch6-discovery/transport"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gofrs/uuid"
)

func main() {
	// 命令行讀取參數，沒有就使用預設值

	var (
		// 服務資訊
		serviceName = flag.String("service.name", "HelloWorld", "service name")
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
	// discoveryClient, err := discover.NewMyDiscoverClient(*consulHost, *consulPort)
	discoveryClient, err := discover.NewKitDiscoverClient(*consulHost, *consulPort)

	// 獲取服務失敗
	if err != nil {
		config.Logger.Println("Get Consul Client failed")
		os.Exit(-1)
	}

	// 初始化 Service
	var svc = service.NewDiscoveryServiceImpl(discoveryClient)

	// SayHello的Endpoint
	sayHelloEndpoint := endpoint.MakeSayHelloEndpoint(svc)
	// Discovery的Endpoint
	discoveryEndpoint := endpoint.MakeDiscoveryEndpoint(svc)
	// HealthCheck的Endpoint
	healthEndpoint := endpoint.MakeHealthCheckEndpoint(svc)

	endpts := endpoint.DiscoveryEndpoints{
		SayHelloEndpoint:    sayHelloEndpoint,
		DiscoveryEndpoint:   discoveryEndpoint,
		HealthCheckEndpoint: healthEndpoint,
	}

	// 創建http.Handler
	r := transport.MakeHttpHandler(ctx, endpts, config.KitLogger)
	// 定義服務實例ID
	_uuid, _ := uuid.NewV4()
	instanceId := *serviceName + "-" + _uuid.String()

	// 啟動 http server
	go func() {
		config.Logger.Println("Http Server start at port:" + strconv.Itoa(*servicePort))
		// 啟動前註冊Service
		if !discoveryClient.Register(*serviceName, instanceId, "/health", *serviceHost, *servicePort, nil, config.Logger) {
			config.Logger.Printf("string-service for service %s failed.", *serviceName)
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
