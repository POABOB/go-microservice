package discover

import (
	"log"
	"strconv"
	"sync"

	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

type KitDiscoverClient struct {
	Host         string // Consul Host
	Port         int    // Consul Port
	client       consul.Client
	config       *api.Config // consul 連線配置
	mutex        sync.Mutex  // 互斥鎖
	instancesMap sync.Map    // 暫存服務實例的資訊
}

func NewKitDiscoverClient(consulHost string, consulPort int) (DiscoveryClient, error) {
	// 1. 設定Config和地址
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulHost + ":" + strconv.Itoa(consulPort)

	// 2. 建立 consul.Client
	apiClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}
	client := consul.NewClient(apiClient)
	return &KitDiscoverClient{
		Host:   consulHost,
		Port:   consulPort,
		config: consulConfig,
		client: client,
	}, err
}

func (consulClient *KitDiscoverClient) Register(serviceName, instanceId, healthCheckUrl string, instanceHost string, instancePort int, meta map[string]string, logger *log.Logger) bool {
	// 服務實例的MetaData，呼叫註冊函數
	err := consulClient.client.Register(&api.AgentServiceRegistration{
		ID:      instanceId,
		Name:    serviceName,
		Address: instanceHost,
		Port:    instancePort,
		Meta:    meta,
		Check: &api.AgentServiceCheck{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           "http://" + instanceHost + ":" + strconv.Itoa(instancePort) + healthCheckUrl,
			Interval:                       "15s",
		},
	})
	if err != nil {
		log.Printf("Register Service: %s Error, instanceId %s \n", serviceName, instanceId)
		return false
	}
	log.Printf("Register Service: %s Success, instanceId %s \n", serviceName, instanceId)
	return true
}

func (consulClient *KitDiscoverClient) DeRegister(instanceId string, logger *log.Logger) bool {
	// 服務實例的MetaData，只需要實例ID，呼叫註銷函數
	err := consulClient.client.Deregister(&api.AgentServiceRegistration{
		ID: instanceId,
	})
	if err != nil {
		log.Printf("Deregister Service Error, instanceId %s \n", instanceId)
		return false
	}
	log.Printf("Deregister Service Success, instanceId %s \n", instanceId)

	return true
}

func (consulClient *KitDiscoverClient) DiscoverServices(serviceName string, logger *log.Logger) []interface{} {
	//  緩存查找服務資訊
	instanceList, ok := consulClient.instancesMap.Load(serviceName)
	if ok {
		return instanceList.([]interface{})
	}

	// 互斥鎖，用途是避免相同名稱服務重新註冊
	consulClient.mutex.Lock()
	defer consulClient.mutex.Unlock()
	// 再次檢查是否有被其他服務註冊
	instanceList, ok = consulClient.instancesMap.Load(serviceName)
	if ok {
		return instanceList.([]interface{})
	} else {
		// 註冊的監控
		go func() {
			// 使用 consul 的 watch 來對服務實例監控
			params := make(map[string]interface{})
			params["type"] = "service"
			params["service"] = serviceName
			plan, _ := watch.Parse(params)
			plan.Handler = func(u uint64, i interface{}) {
				if i == nil {
					return
				}
				v, ok := i.([]*api.ServiceEntry)
				if !ok {
					return // 資料異常，忽略
				}
				// 沒有服務實例在線上，存空的值
				if len(v) == 0 {
					consulClient.instancesMap.Store(serviceName, []interface{}{})
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
		consulClient.instancesMap.Store(serviceName, []interface{}{})
		logger.Printf("Discover Service: %s Error!\n", serviceName)
		return nil
	}

	// 服務實例
	instances := make([]interface{}, len(entries))
	for i := 0; i < len(instances); i++ {
		instances[i] = entries[i].Service
	}

	// 存放在暫存的Map中，避免重複呼叫
	consulClient.instancesMap.Store(serviceName, instances)
	return instances
}
