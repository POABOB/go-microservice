package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/POABOB/go-microservice/common/discover"
	loadbalance "github.com/POABOB/go-microservice/common/load-balance"

	kitlog "github.com/go-kit/log"
)

func main() {

	// 環境變數
	var (
		consulHost = flag.String("consul.host", "127.0.0.1", "consul server ip address")
		consulPort = flag.Int("consul.port", 8500, "consul server port")
	)
	flag.Parse()

	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stderr)
		logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}

	consulClient, err := discover.NewKitHTTPDiscoverClient(*consulHost, *consulPort)

	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	// 反向代理
	proxy := NewHystrixHandler(consulClient, new(loadbalance.RandomLoadBalance), log.New(os.Stderr, "", log.LstdFlags))

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// 監聽
	go func() {
		logger.Log("transport", "HTTP", "addr", "9090")
		errc <- http.ListenAndServe(":9090", proxy)
	}()

	// 執行並等待結束
	logger.Log("exit", <-errc)
}
