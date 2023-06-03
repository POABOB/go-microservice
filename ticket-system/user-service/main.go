package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strconv"

	pb_user "github.com/POABOB/go-microservice/ticket-system/pb/user"
	"github.com/POABOB/go-microservice/ticket-system/pkg/bootstrap"
	conf "github.com/POABOB/go-microservice/ticket-system/pkg/config"
	"github.com/POABOB/go-microservice/ticket-system/pkg/discover"
	"github.com/POABOB/go-microservice/ticket-system/pkg/mysql"
	localconfig "github.com/POABOB/go-microservice/ticket-system/user-service/config"
	"github.com/POABOB/go-microservice/ticket-system/user-service/endpoint"
	"github.com/POABOB/go-microservice/ticket-system/user-service/handlers"
	"github.com/POABOB/go-microservice/ticket-system/user-service/service"
	"github.com/POABOB/go-microservice/ticket-system/user-service/transport"

	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	var (
		servicePort = flag.Int("service.port", bootstrap.HttpConfig.Port, "service port")
		grpcAddr    = flag.Int("grpc", bootstrap.RpcConfig.Port, "gRPC listen address.")
	)

	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	var svc pb_user.UserServer = service.UserService{}
	// 依賴注入 service
	svc = handlers.WrapService(svc)
	// service 依賴注入 endpoints
	endpts := endpoint.NewEndpoints(svc)
	// 依賴注入 service endpoints
	endpts = handlers.WrapEndpoints(endpts)

	// Transport 建立 http.Handler
	r := transport.MakeHttpHandler(ctx, endpts, localconfig.ZipkinTracer, localconfig.Logger)

	// http server
	go func() {
		fmt.Println("Http Server start at port:" + strconv.Itoa(*servicePort))
		// 初始化 Mysql
		mysql.InitMysql(conf.MysqlConfig.Host, conf.MysqlConfig.Port, conf.MysqlConfig.User, conf.MysqlConfig.Pwd, conf.MysqlConfig.Db)
		// 啟動前，註冊服務
		discover.Register()
		handler := r
		errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), handler)
	}()

	// grpc server
	go func() {
		fmt.Println("grpc Server start at port" + ":" + strconv.Itoa(*grpcAddr))
		listener, err := net.Listen("tcp", ":"+strconv.Itoa(*grpcAddr))
		if err != nil {
			errChan <- err
			return
		}

		// 註冊 Trace
		serverTracer := kitzipkin.GRPCServerTrace(localconfig.ZipkinTracer, kitzipkin.Name("grpc-transport"))
		tr := localconfig.ZipkinTracer
		md := metadata.MD{}
		parentSpan := tr.StartSpan("grpc-user-service")

		// 將 zipkin 的 Context 注入 metadata.MD{}
		b3.InjectGRPC(&md)(parentSpan.Context())

		// grpc 傳遞 Context
		ctx := metadata.NewIncomingContext(context.Background(), md)
		// Transport 建立 Handler
		handler := transport.MakeGRPCServer(ctx, endpts, serverTracer)
		// 啟動服務
		gRPCServer := grpc.NewServer()
		pb_user.RegisterUserServer(gRPCServer, handler)
		errChan <- gRPCServer.Serve(listener)
	}()

	// 被關閉
	go handlers.InterruptHandler(errChan)

	error := <-errChan
	// 服務關閉後，服務註銷
	discover.Deregister()
	fmt.Println(error)
}
