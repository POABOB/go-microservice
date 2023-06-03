package discover

import (
	"fmt"
	"log"
	"strconv"

	"github.com/POABOB/go-microservice/ticket-system/pkg/common"

	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

// 初始化
func NewDiscoverConsulClient(consulHost string, consulPort int) *DiscoveryClientInstance {
	// 1. 設定Config和地址
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulHost + ":" + strconv.Itoa(consulPort)

	// 2. 建立 consul.Client
	apiClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil
	}
	client := consul.NewClient(apiClient)
	return &DiscoveryClientInstance{
		Host:   consulHost,
		Port:   consulPort,
		client: client,
		config: consulConfig,
	}
}

// 服務註冊
func (consulClient *DiscoveryClientInstance) Register(serviceName, instanceId, instanceHost, healthCheckUrl string, instancePort, instanceWeight int, meta map[string]string, tags []string, logger *log.Logger) bool {
	// 服務實例的MetaData，呼叫註冊函數
	err := consulClient.client.Register(&api.AgentServiceRegistration{
		ID:      instanceId,
		Name:    serviceName,
		Address: instanceHost,
		Port:    instancePort,
		Meta:    meta,
		Tags:    tags,
		Weights: &api.AgentWeights{
			Passing: instanceWeight,
		},
		Check: &api.AgentServiceCheck{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           fmt.Sprintf("http://%v:%v%v", instanceHost, strconv.Itoa(instancePort), healthCheckUrl),
			Interval:                       "15s",
		},
	})
	if err != nil {
		if logger != nil {
			log.Printf("Register Service: %s Failed, instanceId %s, Service Host:Port %s:%d\n", serviceName, instanceId, instanceHost, instancePort)
		}
		return false
	}
	if logger != nil {
		log.Printf("Register Service: %s Success, instanceId %s, Service Host:Port %s:%d\n", serviceName, instanceId, instanceHost, instancePort)
	}

	return true
}

// 服務註銷
func (consulClient *DiscoveryClientInstance) DeRegister(instanceId string, logger *log.Logger) bool {
	// 服務實例的MetaData，只需要實例ID，呼叫註銷函數
	err := consulClient.client.Deregister(&api.AgentServiceRegistration{
		ID: instanceId,
	})
	if err != nil {
		if logger != nil {
			log.Printf("Deregister Service Failed, instanceId %s\n", instanceId)
		}
		return false
	}
	if logger != nil {
		log.Printf("Deregister Service Success, instanceId %s\n", instanceId)
	}

	return true
}

// 服務發現
func (consulClient *DiscoveryClientInstance) DiscoverServices(serviceName string, logger *log.Logger) []*common.ServiceInstance {
	//  cache 查找服務資訊
	instanceList, ok := consulClient.instancesMap.Load(serviceName)
	if ok {
		return instanceList.([]*common.ServiceInstance)
	}

	// 互斥鎖，用途是避免相同名稱服務重新註冊
	consulClient.mutex.Lock()
	defer consulClient.mutex.Unlock()
	// 再次檢查是否有被其他服務註冊
	instanceList, ok = consulClient.instancesMap.Load(serviceName)
	if ok {
		return instanceList.([]*common.ServiceInstance)
	} else {
		// 註冊的監控
		go func() {
			// 使用 consul 的 watch 來對服務實例監控
			params := make(map[string]interface{})
			params["type"] = "service"
			params["service"] = serviceName
			plan, _ := watch.Parse(params)
			// 只要 serviceName 的服務狀態有改變，就會觸發 Handler
			plan.Handler = func(_ uint64, i interface{}) {
				if i == nil {
					return
				}
				v, ok := i.([]*api.ServiceEntry)
				if !ok {
					return // 資料異常，忽略
				}

				// 沒有服務實例在線上，存空的值
				if len(v) == 0 {
					consulClient.instancesMap.Store(serviceName, []*common.ServiceInstance{})
				}

				var healthServices []interface{}
				for _, service := range v {
					if service.Checks.AggregatedStatus() == api.HealthPassing {
						healthServices = append(healthServices, service.Service)
					}
				}
				consulClient.instancesMap.Store(serviceName, healthServices)
			}
			defer plan.Stop()
			plan.Run(consulClient.config.Address)
		}()
	}

	// 使用服務名稱來獲取本服務註冊的資訊
	entries, _, err := consulClient.client.Service(serviceName, "", false, nil)
	if err != nil {
		consulClient.instancesMap.Store(serviceName, []*common.ServiceInstance{})
		if logger != nil {
			logger.Printf("Discover Service: %s Error!\n", serviceName)
		}
		return nil
	}

	// 服務實例
	instances := make([]*common.ServiceInstance, len(entries))
	for i := 0; i < len(instances); i++ {
		instances[i] = newServiceInstance(entries[i].Service)
	}

	// 存放在暫存的Map中，避免重複呼叫
	consulClient.instancesMap.Store(serviceName, instances)
	return instances
}

// 初始化 ServiceInstance
func newServiceInstance(service *api.AgentService) *common.ServiceInstance {
	// rpc 如果沒有設定，那就是該服務的 Port - 1
	rpcPort := service.Port - 1
	if service.Meta != nil {
		if rpcPortString, ok := service.Meta["rpcPort"]; ok {
			rpcPort, _ = strconv.Atoi(rpcPortString)
		}
	}
	return &common.ServiceInstance{
		Host:     service.Address,
		Port:     service.Port,
		GrpcPort: rpcPort,
		Weight:   service.Weights.Passing,
	}

}
