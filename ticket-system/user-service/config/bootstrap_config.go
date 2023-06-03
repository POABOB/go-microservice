package config

import (
	"fmt"
	"log"

	"github.com/POABOB/go-microservice/ticket-system/pkg/bootstrap"

	"github.com/spf13/viper"
)

func Init() {
	// 設定有預設值的 Env
	viper.AutomaticEnv()
	// 設定 Bootstrap 文件
	initBootstrapConfig()

	// 讀取配置
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("err:%s\n", err)
	}

	// 解析 HttpConfig
	if err := subParse("http", &bootstrap.HttpConfig); err != nil {
		log.Fatal("Fail to parse Http config", err)
	}

	// 解析 DiscoverConfig
	if err := subParse("discover", &bootstrap.DiscoverConfig); err != nil {
		log.Fatal("Fail to parse Discover config", err)
	}

	// 解析 ConfigServerConfig
	if err := subParse("config", &bootstrap.ConfigServerConfig); err != nil {
		log.Fatal("Fail to parse config server", err)
	}

	// 解析 RpcConfig
	if err := subParse("rpc", &bootstrap.RpcConfig); err != nil {
		log.Fatal("Fail to parse rpc server", err)
	}
}

func initBootstrapConfig() {
	// 配置名稱 'bootstrap'
	viper.SetConfigName("bootstrap")
	// 路徑
	viper.AddConfigPath("./user-service")
	//windows 環境下為 %GOPATH，linux 環境下為 $GOPATH
	// viper.AddConfigPath("$GOPATH/src/")
	// 配置文件類型
	viper.SetConfigType("yaml")
}

func subParse(key string, value interface{}) error {
	log.Printf("配置文件的前綴為：%v", key)
	sub := viper.Sub(key)
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}
