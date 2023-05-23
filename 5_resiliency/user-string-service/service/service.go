package service

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"use-string-service/config"

	loadbalance "use-string-service/load-balance"

	"github.com/POABOB/go-microservice/common/discover"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/hashicorp/consul/api"
)

// Service constants
const (
	StringServiceCommandName = "UseStringService"
	StringService            = "StringService"
)

// Service Interface
type Service interface {
	// 远程调用 string-service 服务
	UseStringService(operationType, a, b string) (string, error)

	// 健康检查
	HealthCheck() bool
}

// Class StringService
type UseStringService struct {
	// 服務發現 Client
	discoveryClient discover.DiscoveryClient
	loadbalance     loadbalance.LoadBalance
}

type StringResponse struct {
	Result string `json:"result"`
	Error  error  `json:"error"`
}

func NewUseStringService(client discover.DiscoveryClient, lb loadbalance.LoadBalance) Service {
	hystrix.ConfigureCommand(StringServiceCommandName, hystrix.CommandConfig{
		// 觸發最低請求閥值 5
		RequestVolumeThreshold: 5,
	})
	return &UseStringService{
		discoveryClient: client,
		loadbalance:     lb,
	}
}

// 調用
func (s UseStringService) UseStringService(operationType, a, b string) (string, error) {
	var operationResult string
	var err error

	// 服務發現
	instances := s.discoveryClient.DiscoverServices(StringService, config.Logger)
	instanceList := make([]*api.AgentService, len(instances))
	for i := 0; i < len(instances); i++ {
		instanceList[i] = instances[i].(*api.AgentService)
	}
	// loadbalance 選擇服務
	selectInstance, err := s.loadbalance.SelectService(instanceList)
	if err == nil {
		config.Logger.Printf("current string-service ID is %s and address:port is %s:%s\n", selectInstance.ID, selectInstance.Address, strconv.Itoa(selectInstance.Port))
		requestUrl := url.URL{
			Scheme: "http",
			Host:   selectInstance.Address + ":" + strconv.Itoa(selectInstance.Port),
			Path:   "/op/" + operationType + "/" + a + "/" + b,
		}

		resp, err := http.Post(requestUrl.String(), "", nil)
		if err == nil {
			result := &StringResponse{}
			err = json.NewDecoder(resp.Body).Decode(result)
			if err == nil && result.Error == nil {
				operationResult = result.Result
			}

		}
	}
	return operationResult, err
}

// HealthCheck implement Service method
// 只返回true，暫時不實現
func (s UseStringService) HealthCheck() bool {
	return true
}

// ServiceMiddleware 注入 log 的記錄行為
type ServiceMiddleware func(Service) Service
