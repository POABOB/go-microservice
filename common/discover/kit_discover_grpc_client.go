package discover

import (
	"log"
	"strconv"

	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
)

type KitGRPCDiscoverClient struct {
	*KitDiscoverClient
}

func NewKitGRPCDiscoverClient(consulHost string, consulPort int) (DiscoveryClient, error) {
	// 1. 設定Config和地址
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulHost + ":" + strconv.Itoa(consulPort)

	// 2. 建立 consul.Client
	apiClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}
	client := consul.NewClient(apiClient)
	return &KitGRPCDiscoverClient{
		KitDiscoverClient: &KitDiscoverClient{
			Host:   consulHost,
			Port:   consulPort,
			config: consulConfig,
			client: client,
		},
	}, err
}

func (consulClient *KitGRPCDiscoverClient) Register(serviceName, instanceId, healthCheckUrl string, instanceHost string, instancePort int, meta map[string]string, logger *log.Logger) bool {
	// 服務實例的MetaData，呼叫註冊函數
	err := consulClient.KitDiscoverClient.client.Register(&api.AgentServiceRegistration{
		ID:      instanceId,
		Name:    serviceName,
		Address: instanceHost,
		Port:    instancePort,
		Meta:    meta,
		Check: &api.AgentServiceCheck{
			DeregisterCriticalServiceAfter: "30s",
			GRPC:                           instanceHost + ":" + strconv.Itoa(instancePort) + healthCheckUrl,
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
