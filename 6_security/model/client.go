package model

type ClientDetails struct {
	ClientId                    string   // client 獨立 ID
	ClientSecret                string   // client 的私鑰
	AccessTokenValiditySeconds  int      // 訪問的 Token 有效時間(秒)
	RefreshTokenValiditySeconds int      // 刷新的 Token 有效時間(秒)
	RegisteredRedirectUri       string   // 重新導向的位址，授權碼類型中使用
	AuthorizedGrantTypes        []string // 可以授權的類型
}
