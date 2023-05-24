package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"

	"github.com/POABOB/go-microservice/common/discover"
	loadbalance "github.com/POABOB/go-microservice/common/load-balance"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/hashicorp/consul/api"
)

var (
	ErrNoInstances = errors.New("query service instance error")
)

type HystrixHandler struct {
	// 紀錄服務是否有被配置
	hystrixs      map[string]bool
	hystrixsMutex *sync.Mutex

	discoveryClient discover.DiscoveryClient
	loadbalance     loadbalance.LoadBalance
	logger          *log.Logger
}

func NewHystrixHandler(discoveryClient discover.DiscoveryClient, loadbalance loadbalance.LoadBalance, logger *log.Logger) *HystrixHandler {
	return &HystrixHandler{
		discoveryClient: discoveryClient,
		logger:          logger,
		hystrixs:        make(map[string]bool),
		loadbalance:     loadbalance,
		hystrixsMutex:   &sync.Mutex{},
	}

}

func (hystrixHandler *HystrixHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	// 獲取請求路徑
	reqPath := req.URL.Path
	if reqPath == "" {
		return
	}

	// 使用 '/' 對路徑進行切割，獲取 serviceName
	pathArray := strings.Split(reqPath, "/")
	serviceName := pathArray[1]

	// 如果URL路徑不存在 404
	if serviceName == "" {
		rw.WriteHeader(404)
		return
	}

	// 如果 map 裡面沒有紀錄的話，先鎖起來配置好相關設定
	if _, ok := hystrixHandler.hystrixs[serviceName]; !ok {
		hystrixHandler.hystrixsMutex.Lock()
		if _, ok := hystrixHandler.hystrixs[serviceName]; !ok {
			// 把 serviceName 作為 hystrix 的命令
			hystrix.ConfigureCommand(serviceName, hystrix.CommandConfig{
				// 可以額外進行配置
				Timeout:                5000, // 執行超時的時間(毫秒)
				MaxConcurrentRequests:  8,    // 最大並行量
				SleepWindow:            1000, // 過多久後(毫秒)熔斷器再次進行檢查
				ErrorPercentThreshold:  30,   // 錯誤率 請求數量 >= RequestVolumeThreshold，並且錯誤率到達百分比後就會啟動
				RequestVolumeThreshold: 5,    // 請求閥值(10秒内請求數量)，表示至少有5個請求才進行 ErrorPercentThreshold 錯誤百分比計算
			})
			hystrixHandler.hystrixs[serviceName] = true
		}
		hystrixHandler.hystrixsMutex.Unlock()
	}

	err := hystrix.Do(serviceName, func() error {

		// 使用 consul 來進行服務發現
		instances := hystrixHandler.discoveryClient.DiscoverServices(serviceName, hystrixHandler.logger)
		instanceList := make([]*api.AgentService, len(instances))
		for i := 0; i < len(instances); i++ {
			instanceList[i] = instances[i].(*api.AgentService)
		}
		// 使用演算法來選擇服務實例
		selectInstance, err := hystrixHandler.loadbalance.SelectService(instanceList)

		if err != nil {
			return ErrNoInstances
		}

		// 轉介機制 Director
		director := func(req *http.Request) {

			// 重新分配請求路徑，去除服務名稱
			destPath := strings.Join(pathArray[2:], "/")

			hystrixHandler.logger.Println("service id ", selectInstance.ID)

			// 設置 proxy 資訊
			req.URL.Scheme = "http"
			req.URL.Host = fmt.Sprintf("%s:%d", selectInstance.Address, selectInstance.Port)
			req.URL.Path = "/" + destPath
		}

		var proxyError error

		// 返回代理異常，紀錄 hystrix.Do 執行失敗
		errorHandler := func(_ http.ResponseWriter, _ *http.Request, err error) {
			proxyError = err
		}

		proxy := &httputil.ReverseProxy{
			Director:     director,
			ErrorHandler: errorHandler,
		}

		proxy.ServeHTTP(rw, req)

		// 執行錯誤返回 hystrix
		return proxyError

	}, func(e error) error {
		hystrixHandler.logger.Println("proxy error ", e)
		return errors.New("fallback excute")
	})

	// hystrix.Do 執行異常 返回500和錯誤
	if err != nil {
		rw.WriteHeader(500)
		rw.Write([]byte(err.Error()))
	}

}
