package conf

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/POABOB/go-microservice/ticket-system/pkg/bootstrap"
	"github.com/POABOB/go-microservice/ticket-system/pkg/common"
	consulapi "github.com/hashicorp/consul/api"
	"go.uber.org/zap"

	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	_ "github.com/openzipkin/zipkin-go/reporter/recorder"
	"github.com/spf13/viper"
)

const (
	kConfigType string = "CONFIG_TYPE"
)

var (
	ZipkinTracer *zipkin.Tracer
	Logger       *zap.Logger
)

func initDefault() {
	viper.SetDefault(kConfigType, "yaml")
}

func init() {
	// 設定 Logger
	Logger = common.NewLogger(common.SetAppName("ticket-system"), common.SetDevelopment(true), common.SetLevel(zap.DebugLevel), common.SetErrorFileName("error.log"))

	viper.AutomaticEnv()
	// 初始化預設值
	initDefault()

	// 讀取本地
	// // 設定 Bootstrap 文件
	// initConfig()
	// // 讀取配置
	// if err := viper.ReadInConfig(); err != nil {
	// 	fmt.Printf("err:%s\n", err)
	// }

	// 讀取遠端
	if err := LoadConsulConfig(); err != nil {
		Logger.Error(fmt.Sprintf("Fail to load remote consul config %v", err))
	}

	// 獲取 Mysql Config
	if err := Sub("mysql", &MysqlConfig); err != nil {
		Logger.Error(fmt.Sprintf("Fail to parse mysql %v", err))
	}
	// 獲取 Trace Config
	if err := Sub("trace", &TraceConfig); err != nil {
		Logger.Error(fmt.Sprintf("Fail to parse trace %v", err))
	}

	// 初始化 zipkinUrl
	zipkinUrl := "http://" + TraceConfig.Host + ":" + TraceConfig.Port + TraceConfig.Url
	Logger.Info(fmt.Sprintf("zipkin url: %v", zipkinUrl))
	initTracer(zipkinUrl)
}

// 鏈路追蹤初始化
func initTracer(zipkinURL string) {
	var (
		err           error
		useNoopTracer = zipkinURL == ""
		reporter      = zipkinhttp.NewReporter(zipkinURL)
	)
	// defer reporter.Close()

	zEP, _ := zipkin.NewEndpoint(bootstrap.DiscoverConfig.ServiceName, strconv.Itoa(bootstrap.HttpConfig.Port))
	if ZipkinTracer, err = zipkin.NewTracer(
		reporter, zipkin.WithLocalEndpoint(zEP), zipkin.WithNoopTracer(useNoopTracer),
	); err != nil {
		Logger.Error(fmt.Sprintf("err %v", err))
		os.Exit(1)
	}

	if !useNoopTracer {
		Logger.Info(fmt.Sprintf("tracer: Zipkin, type: Native, URL: %v", zipkinURL))
	}
}

// 讀取 Consul 上的 Config
func LoadConsulConfig() error {
	confAddr := "http://" + bootstrap.DiscoverConfig.Host + ":" + strconv.Itoa(bootstrap.DiscoverConfig.Port)
	config := consulapi.DefaultConfig()
	config.Address = confAddr
	client, err := consulapi.NewClient(config)
	if err != nil {
		Logger.Error(fmt.Sprintf("err %v", err))
		return err
	}

	serviceConfig, _, err := client.KV().Get("/config/"+bootstrap.DiscoverConfig.ServiceName, nil)
	if err != nil {
		Logger.Error(fmt.Sprintf("err %v", err))
		return err
	}

	if err = viper.ReadConfig(bytes.NewBuffer(serviceConfig.Value)); err != nil {
		Logger.Error(fmt.Sprintf("找不到配置 - %v", bootstrap.DiscoverConfig.ServiceName))
		return err
	}

	Logger.Info(fmt.Sprintf("Load config from: %v", confAddr))
	return nil
}

// 讀取對應的 config
func Sub(key string, value interface{}) error {
	Logger.Info(fmt.Sprintf("配置文件的前綴為: %v", key))
	sub := viper.Sub(key)
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}

// 使用本地 config
func initConfig() {
	// 配置名稱 'bootstrap'
	viper.SetConfigName("config")
	// 路徑
	viper.AddConfigPath("./")
	//windows 環境下為 %GOPATH，linux 環境下為 $GOPATH
	// viper.AddConfigPath("$GOPATH/src/")
	// 配置文件類型
	viper.SetConfigType("yaml")
}
