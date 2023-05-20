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

	"string-service/config"
	"string-service/endpoint"
	"string-service/plugins"
	"string-service/service"
	"string-service/transport"

	"github.com/POABOB/go-microservice/common/discover"
	uuid "github.com/satori/go.uuid"
)

func main() {
	// 命令行讀取參數，沒有就使用預設值
	var (
		// 服務資訊
		serviceName = flag.String("service.name", "StringService", "service name")
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

	// 將服務注入 LoggingMiddleware
	var svc service.Service = service.StringService{}
	svc = plugins.LoggingMiddleware(config.KitLogger)(svc)

	// 把服務注入至 Endpoint
	endpts := endpoint.StringEndpoints{
		StringEndpoint:      endpoint.MakeStringEndpoint(svc),      // String Service的 Endpoint
		HealthCheckEndpoint: endpoint.MakeHealthCheckEndpoint(svc), // 健檢的Endpoint
	}

	// 建立 http.Handler
	r := transport.MakeHttpHandler(ctx, endpts, config.KitLogger)
	instanceId := *serviceName + "-" + uuid.NewV4().String()

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
