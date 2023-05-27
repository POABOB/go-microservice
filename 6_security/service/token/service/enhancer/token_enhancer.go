package enhancer

import (
	"security/model"
)

// 針對 token 的操作
type TokenEnhancer interface {
	// 組裝 Token
	Enhance(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error)
	// 獲取 Token 訊息
	Extract(tokenValue string) (*model.OAuth2Token, *model.OAuth2Details, error)
}
