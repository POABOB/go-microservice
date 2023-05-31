package conf

import (
	"sync"

	"github.com/coreos/etcd/clientv3"
	"github.com/go-redis/redis"

	// "ticket-system/sk-core/service/srv_limit"
	"github.com/samuel/go-zookeeper/zk"
)

var (
	Redis       RedisConf
	Etcd        EtcdConf
	SecKill     SecKillConf
	MysqlConfig MysqlConf
	TraceConfig TraceConf
	Zk          ZookeeperConf
)

type ZookeeperConf struct {
	ZkConn        *zk.Conn // Zookeeper 連線實例
	SecProductKey string   // 商品 Key
}

type EtcdConf struct {
	EtcdConn          *clientv3.Client // Etcd 連線實例
	EtcdSecProductKey string           // 商品 Key
	Host              string           // Etcd Host
}

type TraceConf struct {
	Host string
	Port string
	Url  string
}

type MysqlConf struct {
	Host string
	Port string
	User string
	Pwd  string
	Db   string
}

// redis配置
type RedisConf struct {
	RedisConn            *redis.Client // Redis 連線實例
	Proxy2layerQueueName string        // Queue 名稱
	Layer2proxyQueueName string        // Queue 名稱
	IdBlackListHash      string        // 黑名單 ID Hash Table
	IpBlackListHash      string        // 黑名單 IP Hash Table
	IdBlackListQueue     string        // 黑名單 ID Queue
	IpBlackListQueue     string        // 黑名單 IP Queue
	Host                 string        // Host
	Password             string        // Pwd
	Db                   int           // Db
}

type SecKillConf struct {
	RedisConf *RedisConf
	EtcdConf  *EtcdConf

	CookieSecretKey string

	ReferWhiteList []string // 白名單

	AccessLimitConf AccessLimitConf

	RWBlackLock                  sync.RWMutex
	WriteProxy2LayerGoroutineNum int
	ReadProxy2LayerGoroutineNum  int

	IPBlackMap map[string]bool
	IDBlackMap map[int]bool

	SecProductInfoMap map[int]*SecProductInfoConf

	AppWriteToHandleGoroutineNum  int
	AppReadFromHandleGoroutineNum int

	CoreReadRedisGoroutineNum  int
	CoreWriteRedisGoroutineNum int
	CoreHandleGoroutineNum     int

	AppWaitResultTimeout int

	CoreWaitResultTimeout int

	MaxRequestWaitTimeout int

	SendToWriteChanTimeout  int //
	SendToHandleChanTimeout int //
	TokenPassWd             string
}

// 商品資訊
type SecProductInfoConf struct {
	ProductId         int                 `json:"product_id"`           // 商品ID
	StartTime         int64               `json:"start_time"`           // 開始時間
	EndTime           int64               `json:"end_time"`             // 結束時間
	Status            int                 `json:"status"`               // 狀態
	Total             int                 `json:"total"`                // 商品總數
	Left              int                 `json:"left"`                 // 商品餘額
	OnePersonBuyLimit int                 `json:"one_person_buy_limit"` // 單個使用者購買限制
	BuyRate           float64             `json:"buy_rate"`             // 購買頻率限制
	SoldMaxLimit      int                 `json:"sold_max_limit"`       // 最大售出限制
	SecLimit          *srv_limit.SecLimit `json:"sec_limit"`            // 限速控制
}

// 訪問限制
type AccessLimitConf struct {
	IPSecAccessLimit   int // IP 每秒訪問限制
	UserSecAccessLimit int // 使用者 每秒訪問限制
	IPMinAccessLimit   int // IP 每分鐘訪問限制
	UserMinAccessLimit int // 使用者 每分鐘訪問限制
}
