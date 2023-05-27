package store

import (
	"errors"
	"security/model"
)

// 錯誤訊息
var (
	ErrNotSupportOperation = errors.New("no support operation")
)

type TokenStore interface {
	// 儲存 AccessToken
	StoreAccessToken(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details)
	// 刪除 AccessToken
	RemoveAccessToken(tokenValue string)
	// 透過 token，獲取 AccessToken 詳細資訊
	ReadAccessToken(tokenValue string) (*model.OAuth2Token, error)
	// 透過 token，獲取 OAuth2Details 資訊
	ReadOAuth2Details(tokenValue string) (*model.OAuth2Details, error)
	// 透過 OAuth2Details，獲取 AccessToken 資訊
	GetAccessToken(oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error)

	// 儲存 RefreshToken
	StoreRefreshToken(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details)
	// 刪除 RefreshToken
	RemoveRefreshToken(oauth2Token string)
	// 透過 token，獲取 RefreshToken 詳細資訊
	ReadRefreshToken(tokenValue string) (*model.OAuth2Token, error)
	// 透過 token，獲取 OAuth2Details 資訊
	ReadOAuth2DetailsForRefreshToken(tokenValue string) (*model.OAuth2Details, error)
}
