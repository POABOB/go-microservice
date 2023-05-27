package service

import (
	"errors"
	"security/model"
)

// 錯誤訊息
var (
	ErrNotSupportGrantType               = errors.New("grant type is not supported")
	ErrInvalidUsernameAndPasswordRequest = errors.New("invalid username, password")
	ErrInvalidTokenRequest               = errors.New("invalid token")
	ErrExpiredToken                      = errors.New("token is expired")
)

// token 服務
type TokenService interface {
	// 根據 AccessToken 獲取 User、Client 資訊
	GetOAuth2DetailsByAccessToken(tokenValue string) (*model.OAuth2Details, error)
	// 根據 User、Client 資訊，建立 AccessToken
	CreateAccessToken(oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error)
	// 根據 RefreshToken 更新 AccessToken
	RefreshAccessToken(refreshTokenValue string) (*model.OAuth2Token, error)
	// 根據 User、Client 資訊獲取已產生的 token
	GetAccessToken(details *model.OAuth2Details) (*model.OAuth2Token, error)
	// 根據 AccessToken 獲取其結構體
	ReadAccessToken(tokenValue string) (*model.OAuth2Token, error)
}
