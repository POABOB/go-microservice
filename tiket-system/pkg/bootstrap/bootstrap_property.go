package bootstrap

var (
	HttpConfig         HttpConf
	DiscoverConfig     DiscoverConf
	ConfigServerConfig ConfigServerConf
	RpcConfig          RpcConf
)

// Http配置
type HttpConf struct {
	Host string
	Port int
}

// RPC配置
type RpcConf struct {
	Port int
}

// 服務發現與著註冊
type DiscoverConf struct {
	Host        string
	Port        int
	ServiceName string
	Weight      int
	InstanceId  string
}

// 配置中心
type ConfigServerConf struct {
	Id      string
	Profile string
	Label   string
}
