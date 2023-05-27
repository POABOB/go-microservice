package service

import (
	"context"
	"net/http"
	"security/model"

	token_service "security/service/token/service"
	user_service "security/service/user"
)

// 此 Granter Class，主要是負責透過 RefreshToken 刷新 AccessToken
type RefreshTokenGranter struct {
	supportGrantType string
	tokenService     token_service.TokenService
}

func NewRefreshGranter(grantType string, userDetailsService user_service.UserDetailsService, tokenService token_service.TokenService) TokenGranter {
	return &RefreshTokenGranter{
		supportGrantType: grantType,
		tokenService:     tokenService,
	}
}

func (tokenGranter *RefreshTokenGranter) Grant(ctx context.Context, grantType string, client *model.ClientDetails, reader *http.Request) (*model.OAuth2Token, error) {
	if grantType != tokenGranter.supportGrantType {
		return nil, ErrNotSupportGrantType
	}

	// 刷新 token 請求
	refreshTokenValue := reader.URL.Query().Get("refresh_token")
	// 如果沒有傳 token
	if refreshTokenValue == "" {
		return nil, ErrInvalidTokenRequest
	}

	// 給 Service 重整
	return tokenGranter.tokenService.RefreshAccessToken(refreshTokenValue)
}
