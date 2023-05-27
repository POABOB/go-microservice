package store

import (
	"security/model"
	token_enhancer "security/service/token/service/enhancer"
)

// 使用 Jwt 的儲存方式
type JwtTokenStore struct {
	jwtTokenEnhancer *token_enhancer.JwtTokenEnhancer
}

// 實例化
func NewJwtTokenStore(jwtTokenEnhancer *token_enhancer.JwtTokenEnhancer) TokenStore {
	return &JwtTokenStore{
		jwtTokenEnhancer: jwtTokenEnhancer,
	}
}

// Jwt 不用儲存 token
func (tokenStore *JwtTokenStore) StoreAccessToken(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details) {
}

// 獲取 AccessToken
func (tokenStore *JwtTokenStore) ReadAccessToken(tokenValue string) (*model.OAuth2Token, error) {
	oauth2Token, _, err := tokenStore.jwtTokenEnhancer.Extract(tokenValue)
	return oauth2Token, err
}

// 使用 token，獲取 OAuth2Details
func (tokenStore *JwtTokenStore) ReadOAuth2Details(tokenValue string) (*model.OAuth2Details, error) {
	_, oauth2Details, err := tokenStore.jwtTokenEnhancer.Extract(tokenValue)
	return oauth2Details, err
}

// TODO 這邊有問題
// 使用 OAuth2Details，獲取 token
func (tokenStore *JwtTokenStore) GetAccessToken(oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error) {
	return nil, ErrNotSupportOperation
}

// Jwt 不用儲存 token
func (tokenStore *JwtTokenStore) RemoveAccessToken(tokenValue string) {
}

// Jwt 不用儲存 token
func (tokenStore *JwtTokenStore) StoreRefreshToken(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details) {
}

// Jwt 不用儲存 token
func (tokenStore *JwtTokenStore) RemoveRefreshToken(oauth2Token string) {
}

// 獲取 ReadRefreshToken
func (tokenStore *JwtTokenStore) ReadRefreshToken(tokenValue string) (*model.OAuth2Token, error) {
	oauth2Token, _, err := tokenStore.jwtTokenEnhancer.Extract(tokenValue)
	return oauth2Token, err
}

// 使用 token，獲取 OAuth2Details
func (tokenStore *JwtTokenStore) ReadOAuth2DetailsForRefreshToken(tokenValue string) (*model.OAuth2Details, error) {
	_, oauth2Details, err := tokenStore.jwtTokenEnhancer.Extract(tokenValue)
	return oauth2Details, err
}
