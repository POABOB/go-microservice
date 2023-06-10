package route

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"

	"github.com/POABOB/go-microservice/ticket-system/gateway/config"
	"github.com/POABOB/go-microservice/ticket-system/pkg/discover"
	"github.com/POABOB/go-microservice/ticket-system/pkg/loadbalance"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/openzipkin/zipkin-go"
	zipkinhttpsvr "github.com/openzipkin/zipkin-go/middleware/http"
	"go.uber.org/zap"
)

// HystrixRouter hystrix 路由
type HystrixRouter struct {
	svcMap      *sync.Map      // 服務實例，儲存已經通過 hystrix 監控服務列表
	logger      *zap.Logger    // Logger
	fallbackMsg string         // fb 訊息
	tracer      *zipkin.Tracer // 服務追蹤對象
	loadbalance loadbalance.LoadBalance
}

func Routes(zipkinTracer *zipkin.Tracer, fbMsg string, logger *zap.Logger) http.Handler {
	return HystrixRouter{
		svcMap:      &sync.Map{},
		logger:      logger,
		fallbackMsg: fbMsg,
		tracer:      zipkinTracer,
		loadbalance: &loadbalance.RandomLoadBalance{},
	}
}

// 路由前綴處理
func preFilter(r *http.Request) bool {
	// 查詢原始路由，如：/string/calculate/10/5
	reqPath := r.URL.Path
	if reqPath == "" {
		return false
	}

	// 判斷該路由是否需要 Authorization
	res := config.Match(reqPath)
	if res {
		return true
	}

	// // 獲取 Token
	// authToken := r.Header.Get("Authorization")
	// if authToken == "" {
	// 	return false
	// }

	// // OAuth 驗證 TODO
	// oauthClient, _ := client.NewOAuthClient("oauth", nil, nil)
	// resp, remoteErr := oauthClient.CheckToken(context.Background(), nil, &pb.CheckTokenRequest{
	// 	Token: authToken,
	// })

	// if remoteErr != nil || resp == nil {
	// 	return false
	// }
	return true
}

// 路由後綴處理
func postFilter() {
	// for custom filter
}

// HTTP Server
func (router HystrixRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 查詢原始路由，如：/string/calculate/10/5
	reqPath := r.URL.Path
	router.logger.Info(fmt.Sprintf("reqPath: %v", reqPath))

	// 健康檢查直接返回
	if reqPath == "/health" {
		w.WriteHeader(200)
		return
	}

	var err error
	if reqPath == "" || !preFilter(r) {
		err = errors.New("illegal request!")
		w.WriteHeader(403)
		w.Write([]byte(err.Error()))
		return
	}

	//按照分隔符'/'对路径进行分解，获取服务名称serviceName
	// 依據 '/' 對路由進行拆解，獲取 serviceName
	pathArray := strings.Split(reqPath, "/")
	serviceName := pathArray[1]

	// 檢查是否已經加入監控
	if _, ok := router.svcMap.Load(serviceName); !ok {
		// 把serviceName 作為命令對象，設定參數
		hystrix.ConfigureCommand(serviceName, hystrix.CommandConfig{Timeout: 1000})
		router.svcMap.Store(serviceName, serviceName)
	}

	// 執行
	err = hystrix.Do(serviceName, func() (err error) {
		// 調用consul api 查詢 serviceName
		serviceInstance, err := discover.DiscoveryService(serviceName)
		if err != nil {
			return err
		}

		director := func(req *http.Request) {
			// 重新組織路徑，去掉 serviceName
			destPath := strings.Join(pathArray[2:], "/")

			// 隨機選取一個實例
			router.logger.Info(fmt.Sprintf("service %v:%v", serviceInstance.Host, serviceInstance.Port))

			// 設定 proxy 資訊
			req.URL.Scheme = "http"
			req.URL.Host = fmt.Sprintf("%s:%d", serviceInstance.Host, serviceInstance.Port)
			// 避免 call 不到 service
			if serviceInstance.Host == "docker.for.mac.localhost" {
				req.URL.Host = fmt.Sprintf("%s:%d", "localhost", serviceInstance.Port)
			}
			req.URL.Path = "/" + destPath
		}

		var proxyError error = nil
		// 為反向代理增加追蹤邏輯，使用 RoundTrip 代替默認的 Transport
		roundTrip, _ := zipkinhttpsvr.NewTransport(router.tracer, zipkinhttpsvr.TransportTrace(true))

		// 反向代理失敗時，錯誤處理
		errorHandler := func(ew http.ResponseWriter, er *http.Request, err error) {
			proxyError = err
		}

		proxy := &httputil.ReverseProxy{
			Director:     director,
			Transport:    roundTrip,
			ErrorHandler: errorHandler,
		}
		proxy.ServeHTTP(w, r)

		return proxyError

	}, func(err error) error {
		//run 執行失敗時，返回 fallback 資訊
		router.logger.Error(fmt.Sprintf("fallback error description: %v", err.Error()))
		return errors.New(router.fallbackMsg)
	})

	// Do() 執行失敗時，HTTP 返回錯誤
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
}
