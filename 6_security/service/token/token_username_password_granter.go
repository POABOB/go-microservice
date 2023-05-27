package service

import (
	"context"
	"net/http"
	"security/model"

	token_service "security/service/token/service"
	user_service "security/service/user"
)

// 此 Granter Class，主要是負責登入時，產生 AccessToken
type UsernamePasswordTokenGranter struct {
	supportGrantType   string
	userDetailsService user_service.UserDetailsService
	tokenService       token_service.TokenService
}

// 實例化
func NewUsernamePasswordTokenGranter(grantType string, userDetailsService user_service.UserDetailsService, tokenService token_service.TokenService) TokenGranter {
	return &UsernamePasswordTokenGranter{
		supportGrantType:   grantType,
		userDetailsService: userDetailsService,
		tokenService:       tokenService,
	}
}

// 授權
func (tokenGranter *UsernamePasswordTokenGranter) Grant(ctx context.Context, grantType string, client *model.ClientDetails, reader *http.Request) (*model.OAuth2Token, error) {
	// 授權類型匹配
	if grantType != tokenGranter.supportGrantType {
		return nil, ErrNotSupportGrantType
	}

	// 結構體獲取登入資訊
	username := reader.FormValue("username")
	email := reader.FormValue("email")
	password := reader.FormValue("password")

	// 表單驗證
	if (username == "" && email == "") || password == "" {
		return nil, ErrInvalidUsernameAndPasswordRequest
	}

	// TODO 找 MYSQL 並返回 User
	// 驗證 User是否存在
	userDetails, err := tokenGranter.userDetailsService.GetUserDetailByUsernameOrEmail(ctx, username, email, password)

	// 不存在
	if err != nil {
		return nil, ErrInvalidUsernameAndPasswordRequest
	}

	// 根據 User、Client 產生token
	return tokenGranter.tokenService.CreateAccessToken(&model.OAuth2Details{
		Client: client,
		User:   userDetails,
	})
}
