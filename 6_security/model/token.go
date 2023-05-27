package model

import "time"

type OAuth2Token struct {
	RefreshToken *OAuth2Token // 刷新的 token
	TokenType    string       // token 的類型
	TokenValue   string       // token
	ExpiresTime  *time.Time   // 過期時間
}

// 判斷 token 是否過期
func (oauth2Token *OAuth2Token) IsExpired() bool {
	return oauth2Token.ExpiresTime != nil &&
		oauth2Token.ExpiresTime.Before(time.Now())
}

// OAuth2 的結構體包含 Client、User
type OAuth2Details struct {
	Client *ClientDetails
	User   *UserDetails
}
