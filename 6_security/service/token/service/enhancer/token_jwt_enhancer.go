package enhancer

import (
	"security/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// CustomClaims，針對 Jwt 進行封裝
type OAuth2TokenCustomClaims struct {
	UserDetails   model.UserDetails
	ClientDetails model.ClientDetails
	RefreshToken  model.OAuth2Token
	jwt.StandardClaims
}

// 要有私鑰
type JwtTokenEnhancer struct {
	secretKey []byte
}

// 實例化
func NewJwtTokenEnhancer(secretKey string) TokenEnhancer {
	return &JwtTokenEnhancer{
		secretKey: []byte(secretKey),
	}
}

// 獲取 Jwt
func (enhancer *JwtTokenEnhancer) Enhance(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error) {
	return enhancer.sign(oauth2Token, oauth2Details)
}

// 解析 Jwt
func (enhancer *JwtTokenEnhancer) Extract(tokenValue string) (*model.OAuth2Token, *model.OAuth2Details, error) {
	var (
		token *jwt.Token
		err   error
	)
	// 解析 Claims
	if token, err = jwt.ParseWithClaims(tokenValue, &OAuth2TokenCustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return enhancer.secretKey, nil
	}); err != nil {
		return nil, nil, err
	}

	// 型別轉換
	claims := token.Claims.(*OAuth2TokenCustomClaims)
	expiresTime := time.Unix(claims.ExpiresAt, 0)

	return &model.OAuth2Token{
			RefreshToken: &claims.RefreshToken,
			TokenValue:   tokenValue,
			ExpiresTime:  &expiresTime,
		}, &model.OAuth2Details{
			User:   &claims.UserDetails,
			Client: &claims.ClientDetails,
		}, nil
}

// 簽名
func (enhancer *JwtTokenEnhancer) sign(oauth2Token *model.OAuth2Token, oauth2Details *model.OAuth2Details) (*model.OAuth2Token, error) {
	var (
		tokenValue string
		err        error
	)

	expireTime := oauth2Token.ExpiresTime
	clientDetails := *oauth2Details.Client
	userDetails := *oauth2Details.User
	clientDetails.ClientSecret = ""
	userDetails.Password = ""

	// 定義 claims 結構體
	claims := OAuth2TokenCustomClaims{
		UserDetails:   userDetails,
		ClientDetails: clientDetails,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "System",
		},
	}

	// 如果有傳遞 RefreshToken 那就直接賦值
	if oauth2Token.RefreshToken != nil {
		claims.RefreshToken = *oauth2Token.RefreshToken
	}

	// 選擇加密方式
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用私鑰加密
	if tokenValue, err = token.SignedString(enhancer.secretKey); err != nil {
		return nil, err
	}

	oauth2Token.TokenValue = tokenValue
	oauth2Token.TokenType = "jwt"
	return oauth2Token, nil
}
