package client

import (
	"context"
	"ticket-system/pb"
	"ticket-system/pkg/discover"
	"ticket-system/pkg/loadbalance"

	"github.com/opentracing/opentracing-go"
)

type OAuthClient interface {
	/***
	 * 檢驗使用者 Token
	 *
	 * @param ctx		context
	 * @param tracer	鏈路追蹤器
	 * @param request	request 請求
	 *
	 * @return (*pb.CheckTokenResponse, error)
	 **/
	CheckToken(ctx context.Context, tracer opentracing.Tracer, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error)
}

type OAuthClientImpl struct {
	manager     ClientManager           // Client 管理器
	serviceName string                  // 服務名稱
	loadBalance loadbalance.LoadBalance // 負載均衡
	tracer      opentracing.Tracer      // 鏈路追蹤器
}

// 檢驗使用者 Token
func (impl *OAuthClientImpl) CheckToken(ctx context.Context, tracer opentracing.Tracer, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	response := new(pb.CheckTokenResponse)
	if err := impl.manager.DecoratorInvoke("/pb.OAuthService/CheckToken", "token_check", tracer, ctx, request, response); err == nil {
		return response, nil
	} else {
		return nil, err
	}
}
func NewOAuthClient(serviceName string, lb loadbalance.LoadBalance, tracer opentracing.Tracer) (OAuthClient, error) {
	if serviceName == "" {
		serviceName = "oauth"
	}
	if lb == nil {
		lb = defaultLoadBalance
	}

	return &OAuthClientImpl{
		manager: &DefaultClientManager{
			serviceName:     serviceName,
			loadBalance:     lb,
			discoveryClient: discover.ConsulService,
			logger:          discover.Logger,
		},
		serviceName: serviceName,
		loadBalance: lb,
		tracer:      tracer,
	}, nil

}
