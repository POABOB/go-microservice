package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/POABOB/go-microservice/common/config"
	"github.com/POABOB/go-microservice/common/discover"
	loadbalance "github.com/POABOB/go-microservice/common/load-balance"
	kitlog "github.com/go-kit/log"
	"github.com/hashicorp/consul/api"

	"github.com/openzipkin/zipkin-go"
	zipkinhttpsvr "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/openzipkin/zipkin-go/reporter"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

func main() {
	// 	環境變數
	var (
		err        error
		consulHost = flag.String("consul.host", "127.0.0.1", "consul server ip address")
		consulPort = flag.Int("consul.port", 8500, "consul server port")
		zipkinURL  = flag.String("zipkin.url", "http://127.0.0.1:9411/api/v2/spans", "Zipkin server url")
	)
	flag.Parse()

	// 建立 Logger
	var logger kitlog.Logger = config.KitLogger

	// 初始化 Tracer
	var zipkinTracer *zipkin.Tracer
	{
		var (
			hostPort      string            = "127.0.0.1:9091"
			serviceName   string            = "gateway-service"
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

	// 建立 Consul client
	var consulClient discover.DiscoveryClient
	if consulClient, err = discover.NewKitHTTPDiscoverClient(*consulHost, *consulPort); err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	// 建立反向代理
	proxy := NewReverseProxy(consulClient, new(loadbalance.RandomLoadBalance), zipkinTracer, config.Logger)

	tags := map[string]string{
		"component": "gateway-server",
	}

	// 建立 zipkin http middleware
	handler := zipkinhttpsvr.NewServerMiddleware(
		zipkinTracer,
		zipkinhttpsvr.SpanName("gateway"),
		zipkinhttpsvr.TagResponseSize(true),
		zipkinhttpsvr.ServerTags(tags),
	)(proxy)

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// 監聽
	go func() {
		logger.Log("transport", "HTTP", "addr", "9091")
		errc <- http.ListenAndServe(":9091", handler)
	}()

	// 等待結束
	logger.Log("exit", <-errc)
}

// NewReverseProxy 創建反向代理
func NewReverseProxy(client discover.DiscoveryClient, loadbalance loadbalance.LoadBalance, zikkinTracer *zipkin.Tracer, logger *log.Logger) *httputil.ReverseProxy {
	zikkinTracer.StartSpan("1123")
	// 建立 Director
	director := func(req *http.Request) {

		// 獲取原始路祭，如：/string-service/op/10/5
		reqPath := req.URL.Path
		if reqPath == "" {
			return
		}

		// 根據 '/' 分解路徑，獲取 serviceName
		pathArray := strings.Split(reqPath, "/")
		serviceName := pathArray[1]
		if serviceName == "" {
			return
		}

		// 調用 consul api 查找 serviceName 的服務實例
		instances := client.DiscoverServices(serviceName, logger)
		result := make([]*api.AgentService, len(instances))
		for i := 0; i < len(instances); i++ {
			result[i] = instances[i].(*api.AgentService)
		}
		if len(result) == 0 {
			config.KitLogger.Log("ReverseProxy failed", "no such service instance", serviceName)
			return
		}

		// 重新組織路徑，去除 serviceName
		destPath := strings.Join(pathArray[2:], "/")

		// 隨機選取一個服務
		selectInstance, err := loadbalance.SelectService(result)
		if err != nil {
			config.KitLogger.Log("ReverseProxy failed", "query service instance error", serviceName)
			return
		}
		config.KitLogger.Log("service id", selectInstance.ID)

		// 設定 proxy 地址資訊
		req.URL.Scheme = "http"
		req.URL.Host = fmt.Sprintf("%s:%d", selectInstance.Address, selectInstance.Port)
		req.URL.Path = "/" + destPath
	}

	// 反向代理新增追蹤邏輯
	roundTrip, _ := zipkinhttpsvr.NewTransport(zikkinTracer, zipkinhttpsvr.TransportTrace(true))

	return &httputil.ReverseProxy{
		Director:  director,
		Transport: roundTrip,
	}
}
