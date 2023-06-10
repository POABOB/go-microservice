package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/POABOB/go-microservice/ticket-system/gateway/route"
	"github.com/POABOB/go-microservice/ticket-system/pkg/bootstrap"
	conf "github.com/POABOB/go-microservice/ticket-system/pkg/config"
	register "github.com/POABOB/go-microservice/ticket-system/pkg/discover"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/openzipkin/zipkin-go"
	zipkinhttpsvr "github.com/openzipkin/zipkin-go/middleware/http"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

func main() {
	// 建立環境變數
	var (
		zipkinURL = flag.String("zipkin.url", "http://127.0.0.1:9411/api/v2/spans", "Zipkin server url")
	)
	flag.Parse()

	// 建立 Logger
	logger := conf.Logger

	var zipkinTracer *zipkin.Tracer
	{
		var (
			err           error
			useNoopTracer = *zipkinURL == ""
			reporter      = zipkinhttp.NewReporter(*zipkinURL)
		)
		defer reporter.Close()
		zEP, _ := zipkin.NewEndpoint(bootstrap.DiscoverConfig.ServiceName, strconv.Itoa(bootstrap.HttpConfig.Port))
		zipkinTracer, err = zipkin.NewTracer(
			reporter, zipkin.WithLocalEndpoint(zEP), zipkin.WithNoopTracer(useNoopTracer),
		)

		if err != nil {
			logger.Error(fmt.Sprintf("err %v", err))
			os.Exit(1)
		}

		if !useNoopTracer {
			logger.Info(fmt.Sprintf("tracer: Zipkin, type: Native, URL: %v", zipkinURL))
		}
	}

	// 服務註冊
	register.Register()

	// 熔斷路由
	hystrixRouter := route.Routes(zipkinTracer, "Circuit Breaker:Service unavailable", logger)

	handler := zipkinhttpsvr.NewServerMiddleware(
		zipkinTracer,
		zipkinhttpsvr.SpanName(bootstrap.DiscoverConfig.ServiceName),
		zipkinhttpsvr.TagResponseSize(true),
		zipkinhttpsvr.ServerTags(map[string]string{
			"component": "gateway_server",
		}),
	)(hystrixRouter)

	errc := make(chan error)

	// 開啟 hystrix 監控，端口 9010
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go func() {
		errc <- http.ListenAndServe(net.JoinHostPort("", "9010"), hystrixStreamHandler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// 開始監聽
	go func() {
		logger.Info(fmt.Sprintf("transport: HTTP, addr: %v", bootstrap.HttpConfig.Port))
		register.Register()
		errc <- http.ListenAndServe(fmt.Sprintf(":%v", bootstrap.HttpConfig.Port), handler)
	}()

	// 服務執行，等待結束
	error := <-errc
	// 服務結束，註銷服務
	register.Deregister()
	logger.Error(fmt.Sprintf("exit: %v", error))
}
