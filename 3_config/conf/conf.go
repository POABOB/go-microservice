package conf

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	_ "github.com/streadway/amqp"
)

// 基本參數定義
const (
	kAppName       = "APP_NAME"
	kConfigServer  = "CONFIG_SERVER"
	kConfigLabel   = "CONFIG_LABEL"
	kConfigProfile = "CONFIG_PROFILE"
	kConfigType    = "CONFIG_TYPE"
	kAmqpURI       = "AmqpURI"
)

var (
	Resume ResumeConfig
)

type ResumeConfig struct {
	Name string
	Age  int
	Sex  string
}

func init() {
	// viper 設定預設值
	viper.AutomaticEnv()
	// 初始化配置
	initDefault()

	// 監聽變化
	go StartListener(viper.GetString(kAppName), viper.GetString(kAmqpURI), "springCloudBus")

	// 讀取遠端配置
	if err := loadRemoteConfig(); err != nil {
		log.Fatal("Fail to load config", err)
	}

	// 將配置反序列化為 Struct
	if err := sub("resume", &Resume); err != nil {
		log.Fatal("Fail to parse config", err)
	}
}

func initDefault() {
	viper.SetDefault(kAppName, "client-demo")                     // 服務名稱
	viper.SetDefault(kConfigServer, "http://127.0.0.1:8888")      // 配置服務的地址
	viper.SetDefault(kConfigLabel, "main")                        // 分支類型
	viper.SetDefault(kConfigProfile, "dev")                       // 環境類型
	viper.SetDefault(kConfigType, "yaml")                         // 配置檔案類型
	viper.SetDefault(kAmqpURI, "amqp://pass:user@127.0.0.1:5672") // rabbitmq 的地址資訊

}

// 處理刷新事件
func handleRefreshEvent(body []byte, consumerTag string) {
	updateToken := &UpdateToken{}
	err := json.Unmarshal(body, updateToken)
	if err != nil {
		log.Printf("Problem parsing UpdateToken: %v", err.Error())
	} else {
		log.Println(consumerTag, updateToken.DestinationService)
		if strings.Contains(updateToken.DestinationService, consumerTag) {
			log.Println("Reloading Viper config from Spring Cloud Config server")
			loadRemoteConfig()
			log.Println(viper.GetString("resume.name"))
		}
	}
}

// 讀取配置中心的配置
func loadRemoteConfig() (err error) {
	// 地址
	confAddr := fmt.Sprintf("%v/%v/%v-%v.%v",
		viper.Get(kConfigServer), viper.Get(kConfigLabel),
		viper.Get(kAppName), viper.Get(kConfigProfile),
		viper.Get(kConfigType))
	resp, err := http.Get(confAddr)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// 配置檔案類型
	viper.SetConfigType(viper.GetString(kConfigType))
	if err = viper.ReadConfig(resp.Body); err != nil {
		return
	}
	log.Println("Load config from: ", confAddr)
	return
}

// 將 key 解析出 sub-tree，然後反序列化為 Struct
func sub(key string, value interface{}) error {
	log.Printf("配置文件的前綴為：%v", key)
	sub := viper.Sub(key)
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}
