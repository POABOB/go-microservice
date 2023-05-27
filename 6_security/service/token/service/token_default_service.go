package service

import (
	"security/model"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"

	token_enhancer "security/service/token/service/enhancer"
	token_store "security/service/token/service/store"
)

// 預設的 token 服務
type DefaultTokenService struct {
	tokenStore    token_store.TokenStore
	tokenEnhancer token_enhancer.TokenEnhancer
}

// 實例化
func NewTokenService(tokenStore token_store.TokenStore, tokenEnhancer token_enhancer.TokenEnhancer) TokenService {
	return &DefaultTokenService{
		tokenStore:    tokenStore,
		tokenEnhancer: tokenEnhancer,
	}
}

// 根據 User、Client 資訊，建立 AccessToken
func (tokenService *DefaultTokenService) CreateAccessToken(oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error) {
	// 從 Store 中，獲取 GetAccessToken
	existToken, err := tokenService.tokenStore.GetAccessToken(oauth2Details)
	var refreshToken, accessToken *model.OAuth2Token
	if err == nil {
		// token 沒有過期，直接返回
		if !existToken.IsExpired() {
			tokenService.tokenStore.StoreAccessToken(existToken, oauth2Details)
			return existToken, nil
		}

		// token 已經過期，兩個都刪除
		tokenService.tokenStore.RemoveAccessToken(existToken.TokenValue)
		if existToken.RefreshToken != nil {
			refreshToken = existToken.RefreshToken
			tokenService.tokenStore.RemoveRefreshToken(refreshToken.TokenType)
		}
	}

	// 如果 RefreshToken 已經沒了，那就重新建立一個
	if refreshToken == nil || refreshToken.IsExpired() {
		if refreshToken, err = tokenService.createRefreshToken(oauth2Details); err != nil {
			return nil, err
		}
	}

	// 重新產生一個 AccessToken
	if accessToken, err = tokenService.createAccessToken(refreshToken, oauth2Details); err == nil {
		// 儲存 AccessToken
		tokenService.tokenStore.StoreAccessToken(accessToken, oauth2Details)
		tokenService.tokenStore.StoreRefreshToken(refreshToken, oauth2Details)
	}

	return accessToken, err
}

// 私有方法，建立 AccessToken
func (tokenService *DefaultTokenService) createAccessToken(refreshToken *model.OAuth2Token, oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error) {
	// 有效的秒數
	validitySeconds := oauth2Details.Client.AccessTokenValiditySeconds
	// 解析秒數型別
	s, _ := time.ParseDuration(strconv.Itoa(validitySeconds) + "s")
	// 現在時間加上有效秒數
	expiredTime := time.Now().Add(s)
	// 建立 AccessToken
	accessToken := &model.OAuth2Token{
		RefreshToken: refreshToken,
		ExpiresTime:  &expiredTime,
		TokenValue:   uuid.NewV4().String(),
	}

	// 如果有其他 claim 簽發方式，例如 Jwt
	if tokenService.tokenEnhancer != nil {
		return tokenService.tokenEnhancer.Enhance(accessToken, oauth2Details)
	}

	return accessToken, nil
}

// 私有方法，建立 RefreshToken
func (tokenService *DefaultTokenService) createRefreshToken(oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error) {
	// 有效的秒數
	validitySeconds := oauth2Details.Client.RefreshTokenValiditySeconds
	// 解析秒數型別
	s, _ := time.ParseDuration(strconv.Itoa(validitySeconds) + "s")
	// 現在時間加上有效秒數
	expiredTime := time.Now().Add(s)
	// 建立 RefreshToken
	refreshToken := &model.OAuth2Token{
		ExpiresTime: &expiredTime,
		TokenValue:  uuid.NewV4().String(),
	}

	if tokenService.tokenEnhancer != nil {
		return tokenService.tokenEnhancer.Enhance(refreshToken, oauth2Details)
	}

	return refreshToken, nil
}

// 根據 RefreshToken 更新 AccessToken
func (tokenService *DefaultTokenService) RefreshAccessToken(refreshTokenValue string) (*model.OAuth2Token, error) {
	var (
		refreshToken, accessToken, oauth2Token *model.OAuth2Token
		oauth2Details                          *model.OAuth2Details
		err                                    error
	)

	// 從 Store 讀取 RefreshToken
	if refreshToken, err = tokenService.tokenStore.ReadRefreshToken(refreshTokenValue); err != nil {
		return nil, err
	} else if refreshToken.IsExpired() {
		// 過期
		return nil, ErrExpiredToken
	}

	// 用 RefreshToken 讀取 OAuth2Details
	if oauth2Details, err = tokenService.tokenStore.ReadOAuth2DetailsForRefreshToken(refreshTokenValue); err != nil {
		return nil, err
	}

	// 用 OAuth2Details 讀取 OAuth2Token
	if oauth2Token, err = tokenService.tokenStore.GetAccessToken(oauth2Details); err == nil {
		tokenService.tokenStore.RemoveAccessToken(oauth2Token.TokenValue)
	}

	// 刪除原本的 AccessToken 和已經使用過的 RefreshToken
	tokenService.tokenStore.RemoveRefreshToken(refreshTokenValue)

	// 建立 AccessToken
	if refreshToken, err = tokenService.createRefreshToken(oauth2Details); err != nil {
		return nil, err
	}

	// 建立 RefreshToken
	if accessToken, err = tokenService.createAccessToken(refreshToken, oauth2Details); err != nil {
		return nil, err
	}

	// 儲存剛建立的 AccessToken 和 RefreshToken
	tokenService.tokenStore.StoreAccessToken(accessToken, oauth2Details)
	tokenService.tokenStore.StoreRefreshToken(refreshToken, oauth2Details)
	return accessToken, err
}

// 根據 User、Client 資訊獲取已產生的 token
func (tokenService *DefaultTokenService) GetAccessToken(details *model.OAuth2Details) (*model.OAuth2Token, error) {
	return tokenService.tokenStore.GetAccessToken(details)
}

// 根據 AccessToken 獲取其結構體
func (tokenService *DefaultTokenService) ReadAccessToken(tokenValue string) (*model.OAuth2Token, error) {
	return tokenService.tokenStore.ReadAccessToken(tokenValue)
}

// 根據 AccessToken 獲取 User、Client 資訊
func (tokenService *DefaultTokenService) GetOAuth2DetailsByAccessToken(tokenValue string) (*model.OAuth2Details, error) {
	var (
		accessToken *model.OAuth2Token
		err         error
	)

	if accessToken, err = tokenService.tokenStore.ReadAccessToken(tokenValue); err != nil {
		return nil, err
	} else if accessToken.IsExpired() {
		return nil, ErrExpiredToken
	}

	return tokenService.tokenStore.ReadOAuth2Details(tokenValue)
}
