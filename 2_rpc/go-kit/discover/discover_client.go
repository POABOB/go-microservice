package discover

import "log"

type DiscoveryClient interface {
	/***
	 * 服務註冊interface
	 * @param serviceName		服務名稱
	 * @param instanceID		實例 ID
	 * @param instanceHost		實例 Host
	 * @param healthCheckURL	健康檢查 URL
	 * @param instancePort		實例 Port
	 * @param meta				實例 MetaData
	 **/
	Register(serviceName, instanceID, instanceHost, healthCheckURL string, instancePort int, meta map[string]string, logger *log.Logger) bool

	/***
	 * 服務註銷interface
	 * @param instanceID		實例 ID
	 **/
	DeRegister(instanceID string, logger *log.Logger) bool

	/***
	 * 服務發現interface
	 * @param serviceName		服務名稱
	 **/
	DiscoverServices(serviceName string, logger *log.Logger) []interface{}
}
