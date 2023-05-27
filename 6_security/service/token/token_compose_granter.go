package service

import (
	"context"
	"net/http"
	"security/model"
)

// 此 Granter Class，主要是負責根據不同 grantType 來使用不同的授權
type ComposeTokenGranter struct {
	TokenGrantDict map[string]TokenGranter
}

// 實例化
func NewComposeTokenGranter(tokenGrantDict map[string]TokenGranter) TokenGranter {
	return &ComposeTokenGranter{
		TokenGrantDict: tokenGrantDict,
	}
}

// 授權
func (tokenGranter *ComposeTokenGranter) Grant(ctx context.Context, grantType string, client *model.ClientDetails, reader *http.Request) (*model.OAuth2Token, error) {
	// 獲取授權類型
	dispatchGranter := tokenGranter.TokenGrantDict[grantType]
	// 如果不存在，代表不支援
	if dispatchGranter == nil {
		return nil, ErrNotSupportGrantType
	}

	// 返回 TokenGranter
	return dispatchGranter.Grant(ctx, grantType, client, reader)
}
