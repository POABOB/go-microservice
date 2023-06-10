package discover

import (
	"sync"

	"github.com/POABOB/go-microservice/ticket-system/pkg/common"
	"go.uber.org/zap"

	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
)

type DiscoveryClientInstance struct {
	Host         string        // Consul Host
	Port         int           // Consul Port
	client       consul.Client // Concul Client
	config       *api.Config   // consul 連線配置
	mutex        sync.Mutex    // 互斥鎖
	instancesMap sync.Map      // 暫存服務實例的資訊
}

type DiscoveryClient interface {
	/***
	 * 服務註冊interface
	 *
	 * @param serviceName		服務名稱
	 * @param instanceID		實例 ID
	 * @param instanceHost		實例 Host
	 * @param healthCheckURL	健康檢查 URL
	 * @param instancePort		實例 Port
	 * @param instanceWeight	實例權重
	 * @param meta				實例 MetaData
	 * @param tags				實例 tags 標記
	 * @param logger			logger
	 *
	 * @return bool
	 **/
	Register(serviceName, instanceID, instanceHost, healthCheckURL string, instancePort, instanceWeight int, meta map[string]string, tags []string, logger *zap.Logger) bool

	/***
	 * 服務註銷interface
	 *
	 * @param instanceID		實例 ID
	 * @param logger			logger
	 *
	 * @return bool
	 **/
	DeRegister(instanceID string, logger *zap.Logger) bool

	/***
	 * 服務發現interface
	 *
	 * @param serviceName		服務名稱
	 * @param logger			logger
	 *
	 * @return []*common.ServiceInstance
	 **/
	DiscoverServices(serviceName string, logger *zap.Logger) []*common.ServiceInstance
}
