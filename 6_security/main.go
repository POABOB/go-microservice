package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"security/config"
	"security/endpoint"
	"security/model"
	"security/service"
	client_service "security/service/client"
	token_granter "security/service/token"
	token_service "security/service/token/service"
	token_enhancer "security/service/token/service/enhancer"
	token_store "security/service/token/service/store"
	user_service "security/service/user"

	"security/transport"
	"strconv"
	"syscall"

	"github.com/POABOB/go-microservice/common/discover"
	uuid "github.com/satori/go.uuid"
)

func main() {

	var (
		servicePort = flag.Int("service.port", 10098, "service port")
		serviceHost = flag.String("service.host", "127.0.0.1", "service host")
		consulPort  = flag.Int("consul.port", 8500, "consul port")
		consulHost  = flag.String("consul.host", "127.0.0.1", "consul host")
		serviceName = flag.String("service.name", "oauth", "service name")
	)

	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	var discoveryClient discover.DiscoveryClient
	discoveryClient, err := discover.NewKitHTTPDiscoverClient(*consulHost, *consulPort)

	if err != nil {
		config.Logger.Println("Get Consul Client failed")
		os.Exit(-1)

	}

	var tokenEnhancer token_enhancer.TokenEnhancer = token_enhancer.NewJwtTokenEnhancer("secret_sdfshdfsdhfuisdhfsduifhsduifihusdfusidf")
	var tokenStore token_store.TokenStore = token_store.NewJwtTokenStore(tokenEnhancer.(*token_enhancer.JwtTokenEnhancer))
	var tokenService token_service.TokenService = token_service.NewTokenService(tokenStore, tokenEnhancer)

	var userDetailsService user_service.UserDetailsService = user_service.NewInMemoryUserDetailsService(
		[]*model.UserDetails{
			{
				Username:    "simple",
				Email:       "simple@example.com",
				Password:    "123456",
				UserId:      1,
				Authorities: []string{"Simple"},
			},
			{
				Username:    "admin",
				Email:       "admin@example.com",
				Password:    "123456",
				UserId:      1,
				Authorities: []string{"Admin"},
			},
		},
	)
	var clientDetailsService client_service.ClientDetailsService = client_service.NewInMemoryClientDetailService(
		[]*model.ClientDetails{
			{
				ClientId:                    "simple",
				ClientSecret:                "123456",
				AccessTokenValiditySeconds:  1800,
				RefreshTokenValiditySeconds: 18000,
				RegisteredRedirectUri:       "http://127.0.0.1",
				AuthorizedGrantTypes:        []string{"password", "refresh_token"},
			},
		},
	)

	var tokenGranter token_granter.TokenGranter = token_granter.NewComposeTokenGranter(map[string]token_granter.TokenGranter{
		"password":      token_granter.NewUsernamePasswordTokenGranter("password", userDetailsService, tokenService),
		"refresh_token": token_granter.NewRefreshGranter("refresh_token", userDetailsService, tokenService),
	})

	var srv service.Service = service.NewCommonService()

	adminEndpoint := endpoint.MakeAdminEndpoint(srv)
	adminEndpoint = endpoint.MakeOAuth2AuthorizationMiddleware(config.KitLogger)(adminEndpoint)
	adminEndpoint = endpoint.MakeAuthorityAuthorizationMiddleware("Admin", config.KitLogger)(adminEndpoint)

	endpts := endpoint.OAuth2Endpoints{
		TokenEndpoint:       endpoint.MakeClientAuthorizationMiddleware(config.KitLogger)(endpoint.MakeTokenEndpoint(tokenGranter, clientDetailsService)),
		CheckTokenEndpoint:  endpoint.MakeClientAuthorizationMiddleware(config.KitLogger)(endpoint.MakeCheckTokenEndpoint(tokenService)),
		HealthCheckEndpoint: endpoint.MakeHealthCheckEndpoint(srv),
		SimpleEndpoint:      endpoint.MakeOAuth2AuthorizationMiddleware(config.KitLogger)(endpoint.MakeSimpleEndpoint(srv)),
		AdminEndpoint:       adminEndpoint,
	}

	instanceId := *serviceName + "-" + uuid.NewV4().String()

	// http server
	go func() {
		config.Logger.Println("Http Server start at port:" + strconv.Itoa(*servicePort))

		// 服務註冊
		if !discoveryClient.Register(*serviceName, instanceId, *serviceHost, "/health", *servicePort, nil, config.Logger) {
			config.Logger.Printf("use-string-service for service %s failed.", *serviceName)
			// 注册失败，服务启动失败
			os.Exit(-1)
		}

		// HTTP Handler
		handler := transport.MakeHttpHandler(ctx, endpts, tokenService, clientDetailsService, config.KitLogger)
		errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	// 服務註銷
	discoveryClient.DeRegister(instanceId, config.Logger)
	config.Logger.Println(error)
}
