package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-kit/log"
	"github.com/hashicorp/consul/api"
)

func main() {

	// 環境變數
	var (
		consulHost = flag.String("consul.host", "127.0.0.1", "consul server ip address")
		consulPort = flag.String("consul.port", "8500", "consul server port")
	)
	flag.Parse()

	// Logger
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// 創建 consul api client
	consulConfig := api.DefaultConfig()
	consulConfig.Address = "http://" + *consulHost + ":" + *consulPort
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	// 創建反向代理
	proxy := NewReverseProxy(consulClient, logger)

	// 錯誤處理
	errc := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// 開啟服務
	go func() {
		logger.Log("transport", "HTTP", "addr", "9090")
		errc <- http.ListenAndServe(":9090", proxy)
	}()

	// 遇到錯誤結束運行
	logger.Log("exit", <-errc)
}

// NewReverseProxy 反向代理處理方式
func NewReverseProxy(client *api.Client, logger log.Logger) *httputil.ReverseProxy {

	// 轉介機制 Director
	director := func(req *http.Request) {

		// 原始請求路徑
		reqPath := req.URL.Path
		if reqPath == "" {
			return
		}
		// 使用 '/' 對路徑進行切割，獲取 serviceName
		pathArray := strings.Split(reqPath, "/")
		serviceName := pathArray[1]

		// 呼叫 consul api 查詢 serviceName 的服務實例列表
		result, _, err := client.Catalog().Service(serviceName, "", nil)
		if err != nil {
			logger.Log("ReverseProxy failed", "query service instance error", err.Error())
			return
		}

		if len(result) == 0 {
			logger.Log("ReverseProxy failed", "no such service instance", serviceName)
			return
		}

		// 去除 serviceName，重新把路徑處理
		destPath := strings.Join(pathArray[2:], "/")

		// 隨機使用一個服務實例
		tgt := result[rand.Int()%len(result)]
		logger.Log("service id", tgt.ServiceID)

		// 設定代理服務訊息
		req.URL.Scheme = "http"
		req.URL.Host = fmt.Sprintf("%s:%d", tgt.ServiceAddress, tgt.ServicePort)
		req.Host = fmt.Sprintf("%s:%d", tgt.ServiceAddress, tgt.ServicePort)
		req.URL.Path = "/" + destPath
	}
	return &httputil.ReverseProxy{Director: director}

}
