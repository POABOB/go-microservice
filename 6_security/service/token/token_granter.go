package service

import (
	"context"
	"errors"
	"net/http"
	"security/model"
)

// 錯誤訊息
var (
	ErrNotSupportGrantType               = errors.New("grant type is not supported")
	ErrNotSupportOperation               = errors.New("no support operation")
	ErrInvalidUsernameAndPasswordRequest = errors.New("invalid username, password")
	ErrInvalidTokenRequest               = errors.New("invalid token")
	ErrExpiredToken                      = errors.New("token is expired")
)

// token 產生器
type TokenGranter interface {
	Grant(ctx context.Context, grantType string, client *model.ClientDetails, reader *http.Request) (*model.OAuth2Token, error)
}
