package conf

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/POABOB/go-microservice/ticket-system/pkg/bootstrap"
	"github.com/POABOB/go-microservice/ticket-system/pkg/discover"

	"github.com/go-kit/kit/log"
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
	Logger       log.Logger
)

func initDefault() {
	viper.SetDefault(kConfigType, "yaml")
}

func init() {
	// 設定 Logger
	Logger = log.NewLogfmtLogger(os.Stderr)
	Logger = log.With(Logger, "ts", log.DefaultTimestampUTC)
	Logger = log.With(Logger, "caller", log.DefaultCaller)
	viper.AutomaticEnv()
	// 初始化預設值
	initDefault()

	// 獲取遠端 Config
	if err := LoadRemoteConfig(); err != nil {
		Logger.Log("Fail to load remote config", err)
	}

	// 獲取 Mysql Config
	if err := Sub("mysql", &MysqlConfig); err != nil {
		Logger.Log("Fail to parse mysql", err)
	}
	// 獲取 Trace Config
	if err := Sub("trace", &TraceConfig); err != nil {
		Logger.Log("Fail to parse trace", err)
	}

	// 初始化 zipkinUrl
	zipkinUrl := "http://" + TraceConfig.Host + ":" + TraceConfig.Port + TraceConfig.Url
	Logger.Log("zipkin url", zipkinUrl)
	initTracer(zipkinUrl)
}

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
		Logger.Log("err", err)
		os.Exit(1)
	}

	if !useNoopTracer {
		Logger.Log("tracer", "Zipkin", "type", "Native", "URL", zipkinURL)
	}
}
func LoadRemoteConfig() (err error) {
	// 獲取 Config Server
	serviceInstance, err := discover.DiscoveryService(bootstrap.ConfigServerConfig.Id)
	if err != nil {
		return
	}

	configServer := "http://" + serviceInstance.Host + ":" + strconv.Itoa(serviceInstance.Port)
	confAddr := fmt.Sprintf("%v/%v/%v-%v.%v",
		configServer, bootstrap.ConfigServerConfig.Label,
		bootstrap.DiscoverConfig.ServiceName, bootstrap.ConfigServerConfig.Profile,
		viper.Get(kConfigType))
	resp, err := http.Get(confAddr)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// 設定 Config 類型
	viper.SetConfigType(viper.GetString(kConfigType))
	// 讀取 Config
	if err = viper.ReadConfig(resp.Body); err != nil {
		return
	}
	Logger.Log("Load config from: ", confAddr)
	return
}

func Sub(key string, value interface{}) error {
	Logger.Log("配置文件的前綴為：", key)
	sub := viper.Sub(key)
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}
