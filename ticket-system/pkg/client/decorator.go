package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/POABOB/go-microservice/ticket-system/pkg/bootstrap"
	conf "github.com/POABOB/go-microservice/ticket-system/pkg/config"
	"github.com/POABOB/go-microservice/ticket-system/pkg/discover"
	"github.com/POABOB/go-microservice/ticket-system/pkg/loadbalance"
	"go.uber.org/zap"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"

	"github.com/openzipkin/zipkin-go"
	// zipkin "github.com/openzipkin-contrib/zipkin-go-opentracing"

	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"google.golang.org/grpc"
)

var (
	ErrRPCService = errors.New("no rpc service")
)

var DefaultLoadBalance loadbalance.LoadBalance = &loadbalance.RandomLoadBalance{}

type ClientManager interface {
	DecoratorInvoke(path string, hystrixName string, tracer opentracing.Tracer,
		ctx context.Context, inputVal interface{}, outVal interface{}) (err error)
}

type DefaultClientManager struct {
	serviceName     string
	logger          *zap.Logger
	discoveryClient discover.DiscoveryClient
	loadBalance     loadbalance.LoadBalance
	after           []InvokerAfterFunc
	before          []InvokerBeforeFunc
}

type InvokerAfterFunc func() (err error)

type InvokerBeforeFunc func() (err error)

func NewDefaultClientManager(serviceName string, lb loadbalance.LoadBalance) *DefaultClientManager {
	return &DefaultClientManager{
		serviceName:     serviceName,
		logger:          discover.Logger,
		discoveryClient: discover.ConsulService,
		loadBalance:     lb,
	}
}

// 服務調用的裝飾器
func (manager DefaultClientManager) DecoratorInvoke(path string, hystrixName string,
	tracer opentracing.Tracer, ctx context.Context, inputVal interface{}, outVal interface{}) (err error) {
	// Client 發起前的 callback
	for _, fn := range manager.before {
		if err = fn(); err != nil {
			return err
		}
	}

	// 執行 hystrix 熔斷機制
	if err = hystrix.Do(hystrixName, func() error {
		// 服務發現
		instances := manager.discoveryClient.DiscoverServices(manager.serviceName, manager.logger)

		// 負載均衡選擇服務
		if instance, err := manager.loadBalance.SelectService(instances); err == nil {

			// GrpcPort 設定檢查
			if instance.GrpcPort > 0 {
				// Grpc 調用，其中使用 鏈路追蹤 來記錄服務，還有紀錄 Payload，最後就是設定 Timeout
				if conn, err := grpc.Dial(instance.Host+":"+strconv.Itoa(instance.GrpcPort), grpc.WithInsecure(),
					grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(genTracer(tracer), otgrpc.LogPayloads())), grpc.WithTimeout(1*time.Second)); err == nil {
					if err = conn.Invoke(ctx, path, inputVal, outVal); err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				return ErrRPCService
			}
		} else {
			return err
		}
		return nil
	}, func(e error) error {
		return e
	}); err != nil {
		return err
	}

	// Client 發起後的 callback
	for _, fn := range manager.after {
		if err = fn(); err != nil {
			return err
		}
	}
	return nil
}

// 獲取 opentracing.Tracer 實例
func genTracer(tracer opentracing.Tracer) opentracing.Tracer {
	if tracer != nil {
		return tracer
	}

	zipkinUrl := "http://" + conf.TraceConfig.Host + ":" + conf.TraceConfig.Port + conf.TraceConfig.Url
	zipkinRecorder := bootstrap.HttpConfig.Host + ":" + strconv.Itoa(bootstrap.HttpConfig.Port)

	// 設定收集的服務
	reporter := zipkinhttp.NewReporter(zipkinUrl)
	defer reporter.Close()

	// 設定 Endpoint
	endpoint, err := zipkin.NewEndpoint(bootstrap.DiscoverConfig.ServiceName, zipkinRecorder)
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// 設定配置，建立 tracer
	native, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint), zipkin.WithSharedSpans(true), zipkin.WithTraceID128Bit(true))
	if err != nil {
		fmt.Printf("unable to create Zipkin tracer: %+v\n", err)
		os.Exit(-1)
	}

	// 使用 zipkin-go-opentracing 包覆我們的 tracer
	res := zipkinot.Wrap(native)

	// // 將我們的 tracer 設定為單例模式
	// opentracing.SetGlobalTracer(tracer)
	return res
}
