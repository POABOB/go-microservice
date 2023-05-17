package discover

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// 服務實例結構體
type InstanceInfo struct {
	ID                string                     `json:"ID"`                // 服務實例ID
	Service           string                     `json:"Service,omitempty"` // 服務發現的服務名稱
	Name              string                     `json:"Name"`              // 服務名稱
	Tags              []string                   `json:"Tags,omitempty"`    // 標籤，用於分類過濾
	Address           string                     `json:"Address"`           // 服務實例Host
	Port              int                        `json:"Port"`              // 服務實例Port
	Meta              map[string]string          `json:"Meta,omitempty"`    // MetaData
	EnableTagOverride bool                       `json:"EnableTagOverride"` // 是否允許標籤覆蓋
	Check             `json:"Check,omitempty"`   // 健康檢查配置
	Weights           `json:"Weights,omitempty"` // 權重
}

type Check struct {
	DeregisterCriticalServiceAfter string   `json:"DeregisterCriticalServiceAfter"` // 多久後註銷服務
	Args                           []string `json:"Args,omitempty"`                 // 請求參數
	HTTP                           string   `json:"HTTP"`                           // 健康檢查地址
	Interval                       string   `json:"Interval,omitempty"`             // Consul  主動檢查間隔
	TTL                            string   `json:"TTL,omitempty"`                  // 服務實例 主動維持心跳間隔，與Interval擇一
}

type Weights struct {
	Passing int `json:"Passing"`
	Warning int `json:"Warning"`
}

type MyDiscoverClient struct {
	Host string // Consul 的 Host
	Port int    // Consul 的 Port
}

func NewMyDiscoverClient(consulHost string, consulPort int) (DiscoveryClient, error) {
	return &MyDiscoverClient{
		Host: consulHost,
		Port: consulPort,
	}, nil
}

func (consulClient *MyDiscoverClient) Register(serviceName, instanceId, healthCheckUrl string, instanceHost string, instancePort int, meta map[string]string, logger *log.Logger) bool {
	// 1. 封裝服務實例的MetaData
	instanceInfo := &InstanceInfo{
		ID:                instanceId,
		Name:              serviceName,
		Address:           instanceHost,
		Port:              instancePort,
		Meta:              meta,
		EnableTagOverride: false,
		Check: Check{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           "http://" + instanceHost + ":" + strconv.Itoa(instancePort) + healthCheckUrl,
			Interval:                       "15s",
		},
		Weights: Weights{
			Passing: 10,
			Warning: 1,
		},
	}

	byteData, _ := json.Marshal(instanceInfo)

	// 2. 向 Consul 發送註冊服務請求
	req, err := http.NewRequest("PUT",
		"http://"+consulClient.Host+":"+strconv.Itoa(consulClient.Port)+"/v1/agent/service/register",
		bytes.NewReader(byteData))

	if err == nil {
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
		client := http.Client{}
		resp, err := client.Do(req)

		// 3. 檢查結果
		if err != nil {
			log.Println("Register Service Error!")
		} else {
			resp.Body.Close()
			if resp.StatusCode == 200 {
				log.Println("Register Service Success!")
				return true
			} else {
				log.Println("Register Service Error!")
			}
		}
	}
	return false
}

func (consulClient *MyDiscoverClient) DeRegister(instanceId string, logger *log.Logger) bool {
	// 1.发送注销请求
	req, _ := http.NewRequest("PUT",
		"http://"+consulClient.Host+":"+strconv.Itoa(consulClient.Port)+"/v1/agent/service/deregister/"+instanceId, nil)
	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Deregister Service Error!")
	} else {
		resp.Body.Close()
		if resp.StatusCode == 200 {
			log.Println("Deregister Service Success!")
			return true
		} else {
			log.Println("Deregister Service Error!")
		}
	}
	return false
}

func (consulClient *MyDiscoverClient) DiscoverServices(serviceName string, logger *log.Logger) []interface{} {
	// 1. 從 Consul 中獲取服務實例列表
	req, _ := http.NewRequest("GET",
		"http://"+consulClient.Host+":"+strconv.Itoa(consulClient.Port)+"/v1/health/service/"+serviceName, nil)
	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Discover Service Error!")
	} else if resp.StatusCode == 200 {

		var serviceList []struct {
			Service InstanceInfo `json:"Service"`
		}
		err = json.NewDecoder(resp.Body).Decode(&serviceList)
		resp.Body.Close()
		if err == nil {
			instances := make([]interface{}, len(serviceList))
			for i := 0; i < len(instances); i++ {
				instances[i] = serviceList[i].Service
			}
			return instances
		}
	}
	return nil
}
