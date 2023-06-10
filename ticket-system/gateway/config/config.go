package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/POABOB/go-microservice/ticket-system/pkg/bootstrap"
	"github.com/POABOB/go-microservice/ticket-system/pkg/common"
	conf "github.com/POABOB/go-microservice/ticket-system/pkg/config"

	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	_ "github.com/openzipkin/zipkin-go/reporter/recorder"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	kConfigType = "CONFIG_TYPE"
)

var ZipkinTracer *zipkin.Tracer
var (
	Logger *zap.Logger
	Level  zapcore.Level
	Dev    bool
)

func init() {
	Level = zap.InfoLevel
	if bootstrap.ConfigServerConfig.Profile == "dev" {
		Dev = true
		Level = zap.DebugLevel
	}
	// 設定 Logger
	Logger = common.NewLogger(common.SetAppName(bootstrap.DiscoverConfig.ServiceName), common.SetDevelopment(Dev), common.SetLevel(Level), common.SetErrorFileName("error.log"))

	viper.AutomaticEnv()
	initDefault()

	if err := conf.LoadConsulConfig(); err != nil {
		Logger.Error(fmt.Sprintf("Fail to load remote consul config %v", err))
	}

	// 獲取 Mysql Config
	if err := conf.Sub("mysql", &conf.MysqlConfig); err != nil {
		Logger.Error(fmt.Sprintf("Fail to parse mysql %v", err))
	}
	// 獲取 Trace Config
	if err := conf.Sub("trace", &conf.TraceConfig); err != nil {
		Logger.Error(fmt.Sprintf("Fail to parse trace %v", err))
	}

	zipkinUrl := "http://" + conf.TraceConfig.Host + ":" + conf.TraceConfig.Port + conf.TraceConfig.Url
	Logger.Info(fmt.Sprintf("zipkin url: %v", zipkinUrl))
	initTracer(zipkinUrl)
}

func initDefault() {
	viper.SetDefault(kConfigType, "yaml")
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
		Logger.Error(fmt.Sprintf("err %v", err))
		os.Exit(1)
	}

	if !useNoopTracer {
		Logger.Info(fmt.Sprintf("tracer: Zipkin, type: Native, URL: %v", zipkinURL))
	}
}
