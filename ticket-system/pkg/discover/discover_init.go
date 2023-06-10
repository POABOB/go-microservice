package discover

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/POABOB/go-microservice/ticket-system/pkg/bootstrap"
	"github.com/POABOB/go-microservice/ticket-system/pkg/common"
	"github.com/POABOB/go-microservice/ticket-system/pkg/loadbalance"
	"github.com/POABOB/go-microservice/ticket-system/user-service/config"
	"go.uber.org/zap"

	uuid "github.com/satori/go.uuid"
)

var (
	ConsulService        DiscoveryClient                                             // DiscoveryClient 實例
	LoadBalance          loadbalance.LoadBalance                                     // LoadBalance 實例
	Logger               *zap.Logger                                                 // Looger 實例
	ErrNoInstanceExisted error                   = errors.New("no available client") // 錯誤：沒有實例
)

// 建構子
func init() {
	// Discover 初始化
	ConsulService = NewDiscoverConsulClient(bootstrap.DiscoverConfig.Host, bootstrap.DiscoverConfig.Port)
	LoadBalance = new(loadbalance.WeightRoundRobinLoadBalance)
	Logger = config.Logger
}

// 健檢
func CheckHealth(writer http.ResponseWriter, _ *http.Request) {
	Logger.Info("Health check!")
	_, err := fmt.Fprintln(writer, "Server is OK!")
	if err != nil {
		Logger.Error(fmt.Sprintf("%v", err))
	}
}

// 服務發現，使用負載均衡封裝
func DiscoveryService(serviceName string) (*common.ServiceInstance, error) {
	instances := ConsulService.DiscoverServices(serviceName, Logger)

	if len(instances) < 1 {
		Logger.Error(fmt.Sprintf("no available client for %s.", serviceName))
		return nil, ErrNoInstanceExisted
	}
	return LoadBalance.SelectService(instances)
}

// 服務註冊
func Register() {
	// Consul 獲取失敗，關閉
	if ConsulService == nil {
		panic(0)
	}

	// 如果沒有 InstanceId，使用 UUID 獲取 InstanceId
	instanceId := bootstrap.DiscoverConfig.InstanceId
	if instanceId == "" {
		instanceId = bootstrap.DiscoverConfig.ServiceName + uuid.NewV4().String()
	}

	if !ConsulService.Register(bootstrap.DiscoverConfig.ServiceName, instanceId, bootstrap.HttpConfig.Host, "/health",
		bootstrap.HttpConfig.Port, bootstrap.DiscoverConfig.Weight,
		map[string]string{
			"rpcPort": strconv.Itoa(bootstrap.RpcConfig.Port),
		}, nil, Logger) {
		Logger.Error(fmt.Sprintf("register service %s failed.", bootstrap.DiscoverConfig.ServiceName))
		// 註冊失敗，關閉
		panic(0)
	}
	Logger.Info(fmt.Sprintf("%s-service for service %s success.", bootstrap.DiscoverConfig.ServiceName, bootstrap.DiscoverConfig.ServiceName))
}

// 服務註銷
func Deregister() {
	// Consul 獲取失敗，關閉
	if ConsulService == nil {
		panic(0)
	}

	// 如果沒有 InstanceId，使用 UUID 獲取 InstanceId
	instanceId := bootstrap.DiscoverConfig.InstanceId
	if instanceId == "" {
		instanceId = bootstrap.DiscoverConfig.ServiceName + uuid.NewV4().String()
	}

	if !ConsulService.DeRegister(instanceId, Logger) {
		Logger.Error(fmt.Sprintf("deregister for service %s failed.", bootstrap.DiscoverConfig.ServiceName))
		panic(0)
	}
}
