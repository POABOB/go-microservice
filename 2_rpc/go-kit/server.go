package main

import (
	"context"
	"flag"
	"fmt"
	pb "go-kit/pb"
	service "go-kit/string-service"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"go-kit/config"
	"go-kit/discover"

	"github.com/go-kit/kit/log"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc"
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

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	discoveryClient, err := discover.NewKitGRPCDiscoverClient(*consulHost, *consulPort)
	// 獲取服務失敗
	if err != nil {
		logger.Log()
		config.Logger.Println("Get Consul Client failed")
		os.Exit(-1)
	}

	var svc service.Service = service.StringService{}

	// add logging middleware
	svc = service.LoggingMiddleware(logger)(svc)

	//把算术运算Endpoint和健康检查Endpoint封装至StringEndpoints
	endpts := service.StringEndpoints{
		StringEndpoint:      service.MakeStringEndpoint(svc),
		HealthCheckEndpoint: service.MakeHealthCheckEndpoint(svc),
	}

	handler := service.NewStringServer(ctx, endpts)

	ls, _ := net.Listen("tcp", fmt.Sprintf("%s:%s", *serviceHost, strconv.Itoa(*servicePort)))
	gRPCServer := grpc.NewServer()
	pb.RegisterStringServiceServer(gRPCServer, handler)

	uu, _ := uuid.NewV4()
	instanceId := *serviceName + "-" + uu.String()

	// 啟動 http server
	go func() {
		config.Logger.Println("Http Server start at port:" + strconv.Itoa(*servicePort))
		// 啟動前註冊Service
		if !discoveryClient.Register(*serviceName, instanceId, *serviceHost, *serviceName, *servicePort, nil, config.Logger) {
			config.Logger.Printf("string-service for service %s failed.", *serviceName)
			// 註冊失敗，啟動失敗
			os.Exit(-1)
		}
		errChan <- gRPCServer.Serve(ls)
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
